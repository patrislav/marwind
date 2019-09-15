package x11

import (
	"github.com/BurntSushi/xgb/xproto"
)

// WarpPointer moves the pointer to an x, y point on the screen
func WarpPointer(x, y uint32) error {
	return xproto.WarpPointerChecked(
		X, xproto.WindowNone, Screen.Root,
		0, 0, 0, 0,
		int16(x), int16(y),
	).Check()
}
