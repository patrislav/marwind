package container

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind-wm/x11"
)

type WinType uint8

const (
	WinTypeUnknown WinType = iota
	WinTypeNormal
	WinTypeDock
)

type Frame struct {
	col    *Column
	height uint32
	window xproto.Window
	typ    WinType
}

func ManageWindow(xc *xgb.Conn, win xproto.Window) (*Frame, error) {
	typ := getWindowType(xc, win)
	cfgCookie := xproto.ConfigureWindowChecked(xc, win, xproto.ConfigWindowBorderWidth, []uint32{0})
	if err := cfgCookie.Check(); err != nil {
		return nil, err
	}
	evtMask := []uint32{xproto.EventMaskStructureNotify | xproto.EventMaskEnterWindow}
	changeCookie := xproto.ChangeWindowAttributesChecked(xc, win, xproto.CwEventMask, evtMask)
	if err := changeCookie.Check(); err != nil {
		return nil, err
	}
	return &Frame{window: win, typ: typ}, nil
}

func (f *Frame) Height() uint32        { return f.height }
func (f *Frame) Window() xproto.Window { return f.window }
func (f *Frame) Type() WinType         { return f.typ }

func getWindowType(xc *xgb.Conn, win xproto.Window) WinType {
	typeAtom := x11.Atom(xc, "_NET_WM_WINDOW_TYPE")
	dockTypeAtom := x11.Atom(xc, "_NET_WM_WINDOW_TYPE_DOCK")
	normalTypeAtom := x11.Atom(xc, "_NET_WM_WINDOW_TYPE_NORMAL")
	prop, err := xproto.GetProperty(xc, false, win, typeAtom, xproto.GetPropertyTypeAny, 0, 64).Reply()
	if err != nil {
		return WinTypeUnknown
	}
	if prop != nil {
		for v := prop.Value; len(v) >= 4; v = v[4:] {
			switch xproto.Atom(uint32(v[0]) | uint32(v[1])<<8 | uint32(v[2])<<16 | uint32(v[3])<<24) {
			case dockTypeAtom:
				return WinTypeDock
			case normalTypeAtom:
				return WinTypeNormal
			}
		}
	}
	return WinTypeNormal
}
