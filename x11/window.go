package x11

import (
	"fmt"
	"time"

	"github.com/BurntSushi/xgb/xproto"
)

func (xc *Connection) GracefullyDestroyWindow(win xproto.Window) error {
	protos, err := xc.getProps32(win, "WM_PROTOCOLS")
	if err != nil {
		return fmt.Errorf("could not close window: %v", err)
	}
	for _, p := range protos {
		if xproto.Atom(p) == xc.Atom("WM_DELETE_WINDOW") {
			t := time.Now().Unix()
			return xproto.SendEventChecked(
				xc.conn,
				false,
				win,
				xproto.EventMaskNoEvent,
				string(xproto.ClientMessageEvent{
					Format: 32,
					Window: win,
					Type:   xc.Atom("WM_PROTOCOLS"),
					Data: xproto.ClientMessageDataUnionData32New([]uint32{
						uint32(xc.Atom("WM_DELETE_WINDOW")),
						uint32(t),
						0,
						0,
						0,
					}),
				}.Bytes()),
			).Check()
		}
	}
	// The window does not follow ICCCM - just destroy it
	return xproto.DestroyWindowChecked(xc.conn, win).Check()
}

func (xc *Connection) GetRootWindow() xproto.Window {
	return xc.screen.Root
}

func (xc *Connection) CreateWindow(parent xproto.Window, x, y int16, width, height, borderWidth,
	class uint16, valueMask uint32, valueList []uint32) (xproto.Window, error) {

	id, err := xproto.NewWindowId(xc.conn)
	if err != nil {
		return 0, err
	}
	visual := xc.screen.RootVisual
	vdepth := xc.screen.RootDepth
	err = xproto.CreateWindowChecked(xc.conn, vdepth, id, parent, x, y, width, height,
		borderWidth, class, visual, valueMask, valueList).Check()
	if err != nil {
		return 0, fmt.Errorf("could not create window: %s", err)
	}
	return id, nil
}

func (xc *Connection) MapWindow(window xproto.Window) error {
	return xproto.MapWindowChecked(xc.conn, window).Check()
}

func (xc *Connection) UnmapWindow(window xproto.Window) error {
	return xproto.UnmapWindowChecked(xc.conn, window).Check()
}

func (xc *Connection) DestroyWindow(window xproto.Window) error {
	return xproto.DestroyWindowChecked(xc.conn, window).Check()
}

func (xc *Connection) ReparentWindow(window, parent xproto.Window, x, y int16) error {
	return xproto.ReparentWindowChecked(xc.conn, window, parent, x, y).Check()
}
