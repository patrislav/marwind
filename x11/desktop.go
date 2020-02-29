package x11

import (
	"github.com/BurntSushi/xgb/xproto"
)

const (
	leftPtr = 68
)

func (xc *Connection) initDesktop() error {
	cursor, err := xc.createCursor(leftPtr)
	if err != nil {
		return err
	}
	if err := xproto.ChangeWindowAttributesChecked(
		xc.conn,
		xc.screen.Root,
		xproto.CwCursor,
		[]uint32{
			uint32(cursor),
		},
	).Check(); err != nil {
		return err
	}
	return nil
}

func (xc *Connection) createCursor(cursor uint16) (xproto.Cursor, error) {
	fontID, err := xproto.NewFontId(xc.conn)
	if err != nil {
		return 0, err
	}

	cursorID, err := xproto.NewCursorId(xc.conn)
	if err != nil {
		return 0, err
	}

	err = xproto.OpenFontChecked(xc.conn, fontID,
		uint16(len("cursor")), "cursor").Check()
	if err != nil {
		return 0, err
	}

	err = xproto.CreateGlyphCursorChecked(xc.conn, cursorID, fontID, fontID,
		cursor, cursor+1,
		0, 0, 0,
		0xffff, 0xffff, 0xffff).Check()
	if err != nil {
		return 0, err
	}

	err = xproto.CloseFontChecked(xc.conn, fontID).Check()
	if err != nil {
		return 0, err
	}

	return cursorID, nil
}

// WarpPointer moves the pointer to an x, y point on the screen
func (xc *Connection) WarpPointer(x, y int16) error {
	return xproto.WarpPointerChecked(
		xc.conn, xproto.WindowNone, xc.screen.Root,
		0, 0, 0, 0,
		x, y,
	).Check()
}
