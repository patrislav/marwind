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
	fontId, err := xproto.NewFontId(xc.conn)
	if err != nil {
		return 0, err
	}

	cursorId, err := xproto.NewCursorId(xc.conn)
	if err != nil {
		return 0, err
	}

	err = xproto.OpenFontChecked(xc.conn, fontId,
		uint16(len("cursor")), "cursor").Check()
	if err != nil {
		return 0, err
	}

	err = xproto.CreateGlyphCursorChecked(xc.conn, cursorId, fontId, fontId,
		cursor, cursor+1,
		0, 0, 0,
		0xffff, 0xffff, 0xffff).Check()
	if err != nil {
		return 0, err
	}

	err = xproto.CloseFontChecked(xc.conn, fontId).Check()
	if err != nil {
		return 0, err
	}

	return cursorId, nil
}
