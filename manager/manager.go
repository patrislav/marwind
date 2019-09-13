package manager

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind-wm/container"
	"github.com/patrislav/marwind-wm/keysym"
)

// Manager is an instance of the WM
type Manager struct {
	xc        *xgb.Conn
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
	xc, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}
	m.xc = xc

	m.atoms.wmProtocols = m.getAtom("WM_PROTOCOLS")
	m.atoms.wmDeleteWindow = m.getAtom("WM_DELETE_WINDOW")
	m.atoms.wmTakeFocus = m.getAtom("WM_TAKE_FOCUS")

	return m, nil
}

// Init initializes the window manager
func (m *Manager) Init() error {
	if err := xinerama.Init(m.xc); err != nil {
		return err
	}

	conninfo := xproto.Setup(m.xc)
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

	km, err := keysym.LoadKeyMapping(m.xc)
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
	if m.xc != nil {
		m.xc.Close()
	}
}

// Run starts the manager's event loop
func (m *Manager) Run() error {
	for {
		xev, err := m.xc.WaitForEvent()
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
			xproto.SendEventChecked(m.xc, false, e.Window, xproto.EventMaskStructureNotify, string(ev.Bytes()))

		case xproto.DestroyNotifyEvent:
			fmt.Println("DestroyNotifyEvent")
			m.ws.DeleteWindow(e.Window)
			m.renderWorkspace(m.ws)

		case xproto.MapRequestEvent:
			if attrib, err := xproto.GetWindowAttributes(m.xc, e.Window).Reply(); err != nil || !attrib.OverrideRedirect {
				xproto.MapWindowChecked(m.xc, e.Window)
				err := m.ws.AddWindow(m.xc, e.Window)
				if err != nil {
					log.Println("Failed to create a window:", err)
				} else {
					m.renderWorkspace(m.ws)
				}
				m.setFocus(e.Window, xproto.Timestamp(time.Now().Unix()))
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
	return xproto.ChangeWindowAttributesChecked(m.xc, m.xroot.Root, xproto.CwEventMask, evtMask).Check()
}

func (m *Manager) gatherWindows() error {
	tree, err := xproto.QueryTree(m.xc, m.xroot.Root).Reply()
	if err != nil {
		return err
	}
	if tree == nil {
		return errors.New("could not query window tree")
	}
	for _, w := range tree.Children {
		m.ws.AddWindow(m.xc, w)
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
	cookie := xproto.GetProperty(m.xc, false, win, m.atoms.wmProtocols, xproto.GetPropertyTypeAny, 0, 64)
	prop, err := cookie.Reply()
	if err == nil && m.takeFocusProp(prop, win, time) {
		return nil
	}
	_, err = xproto.SetInputFocusChecked(m.xc, xproto.InputFocusPointerRoot, win, time).Reply()
	return err
}

func (m *Manager) takeFocusProp(prop *xproto.GetPropertyReply, win xproto.Window, time xproto.Timestamp) bool {
	for v := prop.Value; len(v) >= 4; v = v[4:] {
		switch xproto.Atom(uint32(v[0]) | uint32(v[1])<<8 | uint32(v[2])<<16 | uint32(v[3])<<24) {
		case m.atoms.wmTakeFocus:
			_ = xproto.SendEventChecked(
				m.xc,
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
				m.xc,
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

func (m *Manager) getAtom(name string) xproto.Atom {
	reply, err := xproto.InternAtom(m.xc, false, uint16(len(name)), name).Reply()
	if err != nil {
		log.Fatal(err)
	}
	if reply == nil {
		return 0
	}
	return reply.Atom
}
