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

func SetActiveWindow(win xproto.Window) error {
	return changeProp32(Screen.Root, "_NET_ACTIVE_WINDOW", xproto.AtomWindow, uint32(win))
}

func setHints() error {
	atoms := make([]uint32, len(ewmhSupported))
	for i, s := range ewmhSupported {
		atoms[i] = uint32(Atom(s))
	}
	return changeProp32(Screen.Root, "_NET_SUPPORTED", xproto.AtomAtom, atoms...)
}
