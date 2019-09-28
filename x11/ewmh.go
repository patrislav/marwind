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
func GetWindowStruts(win xproto.Window) (*Struts, error) {
	// TODO: support _NET_WM_STRUT_PARTIAL as well
	propName := "_NET_WM_STRUT"
	atom := Atom(propName)
	prop, err := xproto.GetProperty(X, false, win, atom, xproto.AtomCardinal, 2, 32).Reply()
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

func SetWMName(name string) error {
	buf := make([]byte, 0)
	buf = append(buf, name...)
	buf = append(buf, 0)
	return changeProp(Screen.Root, 8, "_NET_WM_NAME", xproto.AtomString, buf)
}

func SetActiveWindow(win xproto.Window) error {
	if win == Screen.Root {
		win = 0
	}
	return changeProp32(Screen.Root, "_NET_ACTIVE_WINDOW", xproto.AtomWindow, uint32(win))
}

func SetNumberOfDesktops(num int) error {
	return changeProp32(Screen.Root, "_NET_NUMBER_OF_DESKTOPS", xproto.AtomCardinal, uint32(num))
}

func SetCurrentDesktop(index int) error {
	return changeProp32(Screen.Root, "_NET_CURRENT_DESKTOP", xproto.AtomCardinal, uint32(index))
}

func SetDesktopViewport(num int) error {
	vals := make([]uint32, num*2)
	return changeProp32(Screen.Root, "_NET_DESKTOP_VIEWPORT", xproto.AtomCardinal, vals...)
}

func SetDesktopNames(names []string) error {
	buf := make([]byte, 0)
	for _, name := range names {
		buf = append(buf, name...)
		buf = append(buf, 0)
	}
	return changeProp(Screen.Root, 8, "_NET_DESKTOP_NAMES", Atom("UTF8_STRING"), buf)
}

func SetClientList(windows []xproto.Window) error {
	vals := make([]uint32, len(windows))
	for i, win := range windows {
		vals[i] = uint32(win)
	}
	return changeProp32(Screen.Root, "_NET_CLIENT_LIST", xproto.AtomWindow, vals...)
}

func SetDesktopHints(names []string, index int, windows []xproto.Window) error {
	var err error
	err = SetNumberOfDesktops(len(names))
	if err != nil {
		return err
	}
	err = SetDesktopViewport(len(names))
	if err != nil {
		return err
	}
	err = SetDesktopNames(names)
	if err != nil {
		return err
	}
	err = SetCurrentDesktop(index)
	if err != nil {
		return err
	}
	err = SetClientList(windows)
	if err != nil {
		return err
	}
	return nil
}

func SetWindowDesktop(win xproto.Window, desktop int) error {
	return changeProp32(win, "_NET_WM_DESKTOP", xproto.AtomCardinal, uint32(desktop))
}

func setHints() error {
	atoms := make([]uint32, len(ewmhSupported))
	for i, s := range ewmhSupported {
		atoms[i] = uint32(Atom(s))
	}
	return changeProp32(Screen.Root, "_NET_SUPPORTED", xproto.AtomAtom, atoms...)
}
