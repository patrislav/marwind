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

func createCursor(cursor uint16) (xproto.Cursor, error) {
	fontId, err := xproto.NewFontId(X)
	if err != nil {
		return 0, err
	}

	cursorId, err := xproto.NewCursorId(X)
	if err != nil {
		return 0, err
	}

	err = xproto.OpenFontChecked(X, fontId,
		uint16(len("cursor")), "cursor").Check()
	if err != nil {
		return 0, err
	}

	err = xproto.CreateGlyphCursorChecked(X, cursorId, fontId, fontId,
		cursor, cursor+1,
		0, 0, 0,
		0xffff, 0xffff, 0xffff).Check()
	if err != nil {
		return 0, err
	}

	err = xproto.CloseFontChecked(X, fontId).Check()
	if err != nil {
		return 0, err
	}

	return cursorId, nil
}
