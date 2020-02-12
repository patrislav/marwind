package x11

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"
)

// Struts represents the values of the _NET_WM_STRUT/_NET_WM_STRUT_PARTIAL properties.
// The extended _NET_WM_STRUT_PARTIAL values are ignored - the WM will fill the entire width of the screen instead.
type Struts struct {
	Left, Right, Top, Bottom uint32
}

// GetWindowStruts returns the values of the window's _NET_WM_STRUT(_PARTIAL) property
func (xc *Connection) GetWindowStruts(win xproto.Window) (*Struts, error) {
	// TODO: support _NET_WM_STRUT_PARTIAL as well
	propName := "_NET_WM_STRUT"
	atom := xc.Atom(propName)
	prop, err := xproto.GetProperty(xc.conn, false, win, atom, xproto.AtomCardinal, 2, 32).Reply()
	if err != nil {
		return nil, err
	}
	if prop == nil {
		return nil, fmt.Errorf("xproto.GetProperty returned a nil reply")
	}
	values := make([]uint32, len(prop.Value)/4)
	for v := prop.Value; len(v) >= 4; v = v[4:] {
		values = append(values, uint32(v[0])|uint32(v[1])<<8|uint32(v[2])<<16|uint32(v[3])<<24)
	}
	if len(values) < 4 {
		return nil, fmt.Errorf("not enough values returned by property %s", propName)
	}
	return &Struts{
		Left:   values[0],
		Right:  values[1],
		Top:    values[2],
		Bottom: values[3],
	}, nil
}

func (xc *Connection) SetWMName(name string) error {
	buf := make([]byte, 0)
	buf = append(buf, name...)
	buf = append(buf, 0)
	return xc.changeProp(xc.screen.Root, 8, "_NET_WM_NAME", xproto.AtomString, buf)
}

func (xc *Connection) GetWindowTitle(window xproto.Window) (string, error) {
	reply, err := xc.getProp(window, "_NET_WM_NAME")
	if err != nil {
		return "", err
	}
	return string(reply.Value), nil
}

func (xc *Connection) SetActiveWindow(win xproto.Window) error {
	if win == xc.screen.Root {
		win = 0
	}
	return xc.changeProp32(xc.screen.Root, "_NET_ACTIVE_WINDOW", xproto.AtomWindow, uint32(win))
}

func (xc *Connection) SetNumberOfDesktops(num int) error {
	return xc.changeProp32(xc.screen.Root, "_NET_NUMBER_OF_DESKTOPS", xproto.AtomCardinal, uint32(num))
}

func (xc *Connection) SetCurrentDesktop(index int) error {
	return xc.changeProp32(xc.screen.Root, "_NET_CURRENT_DESKTOP", xproto.AtomCardinal, uint32(index))
}

func (xc *Connection) SetDesktopViewport(num int) error {
	vals := make([]uint32, num*2)
	return xc.changeProp32(xc.screen.Root, "_NET_DESKTOP_VIEWPORT", xproto.AtomCardinal, vals...)
}

func (xc *Connection) SetDesktopNames(names []string) error {
	buf := make([]byte, 0)
	for _, name := range names {
		buf = append(buf, name...)
		buf = append(buf, 0)
	}
	return xc.changeProp(xc.screen.Root, 8, "_NET_DESKTOP_NAMES", xc.Atom("UTF8_STRING"), buf)
}

func (xc *Connection) SetClientList(windows []xproto.Window) error {
	vals := make([]uint32, len(windows))
	for i, win := range windows {
		vals[i] = uint32(win)
	}
	return xc.changeProp32(xc.screen.Root, "_NET_CLIENT_LIST", xproto.AtomWindow, vals...)
}

func (xc *Connection) SetDesktopHints(names []string, index int, windows []xproto.Window) error {
	var err error
	err = xc.SetNumberOfDesktops(len(names))
	if err != nil {
		return err
	}
	err = xc.SetDesktopViewport(len(names))
	if err != nil {
		return err
	}
	err = xc.SetDesktopNames(names)
	if err != nil {
		return err
	}
	err = xc.SetCurrentDesktop(index)
	if err != nil {
		return err
	}
	err = xc.SetClientList(windows)
	if err != nil {
		return err
	}
	return nil
}

func (xc *Connection) SetWindowDesktop(win xproto.Window, desktop int) error {
	return xc.changeProp32(win, "_NET_WM_DESKTOP", xproto.AtomCardinal, uint32(desktop))
}

func (xc *Connection) setHints() error {
	atoms := make([]uint32, len(ewmhSupported))
	for i, s := range ewmhSupported {
		atoms[i] = uint32(xc.Atom(s))
	}
	return xc.changeProp32(xc.screen.Root, "_NET_SUPPORTED", xproto.AtomAtom, atoms...)
}
