package manager

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind-wm/container"
	"github.com/patrislav/marwind-wm/keysym"
	"github.com/patrislav/marwind-wm/x11"
)

// Manager is an instance of the WM
type Manager struct {
	xroot     xproto.ScreenInfo
	keymap    keysym.Keymap
	outputs   []*container.Output
	activeWin xproto.Window
	atoms     struct {
		wmProtocols    xproto.Atom
		wmDeleteWindow xproto.Atom
		wmTakeFocus    xproto.Atom
	}
	actions []Action
	config  Config

	// Temporary property
	ws *container.Workspace
}

// New initializes a Manager and creates an X11 connection
func New(config Config) (*Manager, error) {
	m := &Manager{config: config}
	err := x11.InitConnection()
	if err != nil {
		return nil, err
	}
	m.atoms.wmProtocols = x11.Atom("WM_PROTOCOLS")
	m.atoms.wmDeleteWindow = x11.Atom("WM_DELETE_WINDOW")
	m.atoms.wmTakeFocus = x11.Atom("WM_TAKE_FOCUS")
	return m, nil
}

// Init initializes the window manager
func (m *Manager) Init() error {
	if err := xinerama.Init(x11.X); err != nil {
		return err
	}

	conninfo := xproto.Setup(x11.X)
	if conninfo == nil {
		return errors.New("could not parse X connection info")
	}
	if len(conninfo.Roots) != 1 {
		return errors.New("wrong number of roots, did xinerama initialize properly?")
	}
	m.xroot = conninfo.Roots[0]

	if err := m.becomeWM(); err != nil {
		if _, ok := err.(xproto.AccessError); ok {
			return errors.New("could not become WM, is another WM already running?")
		}
		return err
	}

	km, err := keysym.LoadKeyMapping(x11.X)
	if err != nil {
		log.Fatal(err)
	}
	m.keymap = *km
	m.actions = initActions(m)
	if err := m.grabKeys(); err != nil {
		log.Fatal(err)
	}

	output := container.NewOutput(container.Rect{0, 0, 1366, 768})
	ws := container.NewWorkspace(container.WorkspaceConfig{Gap: m.config.OuterGap})
	output.AddWorkspace(ws)

	m.outputs = append(m.outputs, output)
	m.ws = ws

	m.gatherWindows()

	return nil
}

// Close cleans up the Manager's resources
func (m *Manager) Close() {
	if x11.X != nil {
		x11.X.Close()
	}
}

// Run starts the manager's event loop
func (m *Manager) Run() error {
	for {
		xev, err := x11.X.WaitForEvent()
		if err != nil {
			log.Println(err)
			continue
		}
		switch e := xev.(type) {
		case xproto.KeyPressEvent:
			if err := m.handleKeyPressEvent(e); err != nil {
				return err
			}
		case xproto.ConfigureRequestEvent:
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
			xproto.SendEventChecked(x11.X, false, e.Window, xproto.EventMaskStructureNotify, string(ev.Bytes()))

		case xproto.DestroyNotifyEvent:
			m.deleteWindow(e.Window)

		case xproto.MapRequestEvent:
			if attrib, err := xproto.GetWindowAttributes(x11.X, e.Window).Reply(); err != nil || !attrib.OverrideRedirect {
				xproto.MapWindowChecked(x11.X, e.Window)
				if err := m.addWindow(e.Window); err != nil {
					log.Println("Failed to manage a window:", err)
				}
			}

		case xproto.EnterNotifyEvent:
			m.setFocus(e.Event, e.Time)

		default:
			log.Println(xev)
		}
	}
}

func (m *Manager) becomeWM() error {
	evtMask := []uint32{
		xproto.EventMaskKeyPress |
			xproto.EventMaskKeyRelease |
			xproto.EventMaskButtonPress |
			xproto.EventMaskButtonRelease |
			xproto.EventMaskStructureNotify |
			xproto.EventMaskSubstructureRedirect,
	}
	return xproto.ChangeWindowAttributesChecked(x11.X, m.xroot.Root, xproto.CwEventMask, evtMask).Check()
}

func (m *Manager) addWindow(win xproto.Window) error {
	frame, err := container.ManageWindow(win)
	if err != nil {
		return err
	}
	switch frame.Type() {
	case container.WinTypeNormal:
		m.ws.AddFrame(frame)
		m.renderWorkspace(m.ws)
		m.setFocus(win, xproto.Timestamp(time.Now().Unix()))
	case container.WinTypeDock:
		m.outputs[0].AddDock(frame)
		m.renderOutput(m.outputs[0])
		m.renderWorkspace(m.ws)
	}
	return nil
}

func (m *Manager) deleteWindow(win xproto.Window) error {
	for _, output := range m.outputs {
		if output.DeleteWindow(win) {
			m.renderOutput(output)
			m.renderWorkspace(m.ws)
			return nil
		}
	}
	return fmt.Errorf("could not find window to delete: %v", win)
}

func (m *Manager) gatherWindows() error {
	tree, err := xproto.QueryTree(x11.X, m.xroot.Root).Reply()
	if err != nil {
		return err
	}
	if tree == nil {
		return errors.New("could not query window tree")
	}
	for _, w := range tree.Children {
		m.addWindow(w)
	}
	return nil
}

func (m *Manager) handleKeyPressEvent(e xproto.KeyPressEvent) error {
	sym := m.keymap[e.Detail][0]
	for _, action := range m.actions {
		if sym == action.sym && e.State == uint16(action.modifiers) {
			return action.act()
		}
	}
	return nil
}

func (m *Manager) setFocus(win xproto.Window, time xproto.Timestamp) error {
	m.activeWin = win
	cookie := xproto.GetProperty(x11.X, false, win, m.atoms.wmProtocols, xproto.GetPropertyTypeAny, 0, 64)
	prop, err := cookie.Reply()
	if err == nil && m.takeFocusProp(prop, win, time) {
		return nil
	}
	_, err = xproto.SetInputFocusChecked(x11.X, xproto.InputFocusPointerRoot, win, time).Reply()
	return err
}

func (m *Manager) takeFocusProp(prop *xproto.GetPropertyReply, win xproto.Window, time xproto.Timestamp) bool {
	for v := prop.Value; len(v) >= 4; v = v[4:] {
		switch xproto.Atom(uint32(v[0]) | uint32(v[1])<<8 | uint32(v[2])<<16 | uint32(v[3])<<24) {
		case m.atoms.wmTakeFocus:
			_ = xproto.SendEventChecked(
				x11.X,
				false,
				win,
				xproto.EventMaskNoEvent,
				string(xproto.ClientMessageEvent{
					Format: 32,
					Window: win,
					Type:   m.atoms.wmProtocols,
					Data: xproto.ClientMessageDataUnionData32New([]uint32{
						uint32(m.atoms.wmTakeFocus),
						uint32(time),
						0,
						0,
						0,
					}),
				}.Bytes()),
			).Check()
			return true
		}
	}
	return false
}

func (m *Manager) grabKeys() error {
	for _, action := range m.actions {
		for _, code := range action.codes {
			cookie := xproto.GrabKeyChecked(
				x11.X,
				false,
				m.xroot.Root,
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
