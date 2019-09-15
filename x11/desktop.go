package x11

import (
	"github.com/BurntSushi/xgb/xproto"
)

const (
	leftPtr = 68
)

func initDesktop() error {
	cursor, err := createCursor(leftPtr)
	if err != nil {
		return err
	}
	if err := xproto.ChangeWindowAttributesChecked(
		X,
		Screen.Root,
		xproto.CwCursor,
		[]uint32{
			uint32(cursor),
		},
	).Check(); err != nil {
		return err
	}
	return nil
}
