package client

import (
	"image"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xgraphics"
)

type x11 interface {
	GetRootWindow() xproto.Window
	CreateWindow(
		parent xproto.Window,
		x int16, y int16, width uint16, height uint16,
		borderWidth uint16,
		class uint16, valueMask uint32, valueList []uint32,
	) (xproto.Window, error)

	MapWindow(window xproto.Window) error
	UnmapWindow(window xproto.Window) error
	DestroyWindow(window xproto.Window) error
	ReparentWindow(window, parent xproto.Window, x, y int16) error

	GetWindowTitle(window xproto.Window) (string, error)
	Atom(name string) xproto.Atom

	NewImage(rect image.Rectangle) *xgraphics.Image
}
