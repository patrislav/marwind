package container

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Frame struct {
	col    *Column
	height uint32
	window xproto.Window
}

func ManageWindow(xc *xgb.Conn, win xproto.Window) (*Frame, error) {
	cfgCookie := xproto.ConfigureWindowChecked(xc, win, xproto.ConfigWindowBorderWidth, []uint32{0})
	if err := cfgCookie.Check(); err != nil {
		return nil, err
	}
	evtMask := []uint32{xproto.EventMaskStructureNotify | xproto.EventMaskEnterWindow}
	changeCookie := xproto.ChangeWindowAttributesChecked(xc, win, xproto.CwEventMask, evtMask)
	if err := changeCookie.Check(); err != nil {
		return nil, err
	}
	return &Frame{window: win}, nil
}

func (f *Frame) Height() uint32        { return f.height }
func (f *Frame) Window() xproto.Window { return f.window }
