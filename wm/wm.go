package wm

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind/client"
	"github.com/patrislav/marwind/keysym"
	"github.com/patrislav/marwind/x11"
)

const maxWorkspaces = 10

// WM is a struct representing the Window Manager
type WM struct {
	xc           *x11.Connection
	outputs      []*output
	keymap       keysym.Keymap
	actions      []*action
	config       Config
	workspaces   [maxWorkspaces]*workspace
	activeWin    xproto.Window
	windowConfig *client.Config
}

// New initializes a WM and creates an X11 connection
func New(config Config) (*WM, error) {
	wc := &client.Config{
		BgColor:        config.BorderColor,
		TitlebarHeight: config.TitleBarHeight,
		FontColor:      config.TitleBarFontColorActive,
		FontSize:       config.TitleBarFontSize,
		BorderWidth:    config.BorderWidth,
	}
	xconn, err := x11.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to create WM: %v", err)
	}
	wm := &WM{xc: xconn, config: config, windowConfig: wc}
	return wm, nil
}

// Init initializes the WM
func (wm *WM) Init() error {
	if err := wm.xc.Init(); err != nil {
		return fmt.Errorf("failed to init WM: %v", err)
	}
	if err := wm.becomeWM(); err != nil {
		if _, ok := err.(xproto.AccessError); ok {
			return fmt.Errorf("could not become WM, possibly another WM is already running")
		}
		return fmt.Errorf("could not become WM: %v", err)
	}
	km, err := keysym.LoadKeyMapping(wm.xc.X())
	if err != nil {
		return fmt.Errorf("failed to load key mapping: %v", err)
	}
	wm.keymap = *km
	wm.actions = initActions(wm)
	if err := wm.grabKeys(); err != nil {
		return fmt.Errorf("failed to grab keys: %v", err)
	}

	o := newOutput(wm.xc, client.Geom{
		X: 0, Y: 0,
		W: wm.xc.Screen().WidthInPixels,
		H: wm.xc.Screen().HeightInPixels,
	})
	for i := 0; i < maxWorkspaces; i++ {
		wm.workspaces[i] = newWorkspace(uint8(i), workspaceConfig{gap: wm.config.OuterGap})
	}
	if err := o.addWorkspace(wm.workspaces[0]); err != nil {
		return fmt.Errorf("failed to add workspace to output: %v", err)
	}
	wm.outputs = append(wm.outputs, o)

	if err := wm.xc.SetWMName("Marwind"); err != nil {
		return fmt.Errorf("failed to set WM name: %v", err)
	}
	if err := wm.manageExistingClients(); err != nil {
		return fmt.Errorf("failed to manage existing clients: %v", err)
	}
	return nil
}

// Close cleans up the WM's resources
func (wm *WM) Close() {
	if wm.xc != nil {
		wm.xc.Close()
	}
}

// Run starts the WM's X event loop
func (wm *WM) Run() error {
	if err := wm.updateDesktopHints(); err != nil {
		return err
	}
	handler := eventHandler{wm: wm}
	handler.eventLoop()
	return nil
}

// becomeWM updates the X root window's attributes in an attempt to manage other windows
func (wm *WM) becomeWM() error {
	evtMask := []uint32{
		xproto.EventMaskKeyPress |
			xproto.EventMaskKeyRelease |
			xproto.EventMaskButtonPress |
			xproto.EventMaskButtonRelease |
			xproto.EventMaskPropertyChange |
			xproto.EventMaskFocusChange |
			xproto.EventMaskStructureNotify |
			xproto.EventMaskSubstructureRedirect,
	}
	return xproto.ChangeWindowAttributesChecked(wm.xc.X(), wm.xc.GetRootWindow(), xproto.CwEventMask, evtMask).Check()
}

// grabKeys attempts to get a sole ownership of certain key combinations
func (wm *WM) grabKeys() error {
	for _, action := range wm.actions {
		for _, code := range action.codes {
			cookie := xproto.GrabKeyChecked(
				wm.xc.X(),
				false,
				wm.xc.GetRootWindow(),
				uint16(action.modifiers),
				code,
				xproto.GrabModeAsync,
				xproto.GrabModeAsync,
			)
			if err := cookie.Check(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (wm *WM) findFrame(predicate func(*frame) bool) *frame {
	for _, ws := range wm.workspaces {
		for _, col := range ws.columns {
			for _, f := range col.frames {
				if predicate(f) {
					return f
				}
			}
		}
	}
	for _, o := range wm.outputs {
		for area := range o.dockAreas {
			for _, f := range o.dockAreas[area] {
				if predicate(f) {
					return f
				}
			}
		}
	}
	return nil
}

func (wm *WM) deleteFrame(f *frame) error {
	for _, o := range wm.outputs {
		if o.deleteFrame(f) {
			if err := wm.removeFocus(); err != nil {
				return err
			}
			return wm.renderOutput(o)
		}
	}
	return fmt.Errorf("could not find frame to delete: %v", f)
}

func (wm *WM) handleKeyPressEvent(e xproto.KeyPressEvent) error {
	sym := wm.keymap[e.Detail][0]
	for _, action := range wm.actions {
		if sym == action.sym && e.State == uint16(action.modifiers) {
			return action.act()
		}
	}
	return nil
}

// TODO: avoid updating all hints at once
func (wm *WM) updateDesktopHints() error {
	out := wm.outputs[0]
	wsWins := make([][]xproto.Window, len(out.workspaces))
	names := make([]string, len(out.workspaces))
	current := 0
	for i, ws := range out.workspaces {
		names[i] = fmt.Sprintf("%d", ws.id+1)
		for _, col := range ws.columns {
			for _, f := range col.frames {
				wsWins[i] = append(wsWins[i], f.cli.Window())
			}
		}
		if ws == out.activeWs {
			current = i
			for area := range out.dockAreas {
				for _, f := range out.dockAreas[area] {
					wsWins[i] = append(wsWins[i], f.cli.Window())
				}
			}
		}
	}
	windows := make([]xproto.Window, 0)
	for _, wins := range wsWins {
		windows = append(windows, wins...)
	}
	if err := wm.xc.SetDesktopHints(names, current, windows); err != nil {
		return err
	}
	var err error
	for i, wins := range wsWins {
		for _, win := range wins {
			if e := wm.xc.SetWindowDesktop(win, i); e != nil {
				err = e
			}
		}
	}
	return err
}

func (wm *WM) handleConfigureRequest(e xproto.ConfigureRequestEvent) error {
	f := wm.findFrame(func(frm *frame) bool { return frm.cli.Window() == e.Window })
	if f != nil {
		if err := wm.configureNotify(f); err != nil {
			return fmt.Errorf("failed to send ConfigureNotify event to %d: %v", e.Window, err)
		}
		return nil
	}
	ev := xproto.ConfigureNotifyEvent{
		Event:            e.Window,
		Window:           e.Window,
		AboveSibling:     0,
		X:                e.X,
		Y:                e.Y,
		Width:            e.Width,
		Height:           e.Height,
		BorderWidth:      0,
		OverrideRedirect: false,
	}
	xproto.SendEventChecked(wm.xc.X(), false, e.Window, xproto.EventMaskStructureNotify, string(ev.Bytes()))
	return nil
}

func (wm *WM) manageExistingClients() error {
	tree, err := xproto.QueryTree(wm.xc.X(), wm.xc.GetRootWindow()).Reply()
	if err != nil {
		return err
	}
	for _, win := range tree.Children {
		attrs, err := xproto.GetWindowAttributes(wm.xc.X(), win).Reply()
		if err != nil {
			continue
		}
		if attrs.MapState == xproto.MapStateUnmapped || attrs.OverrideRedirect {
			continue
		}
		if err := wm.manageWindow(win); err != nil {
			log.Println("Failed to manage an existing window:", err)
		}
	}
	if err := wm.updateDesktopHints(); err != nil {
		return err
	}
	return wm.renderOutput(wm.outputs[0])
}
