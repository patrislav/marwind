package x11

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func changeProp32(win xproto.Window, prop string, typ xproto.Atom, data ...uint32) error {
	buf := make([]byte, len(data)*4)
	for i, datum := range data {
		xgb.Put32(buf[(i*4):], datum)
	}
	return changeProp(win, 32, prop, typ, buf)
}

func changeProp(win xproto.Window, format byte, prop string, typ xproto.Atom, data []byte) error {
	propAtom := Atom(prop)
	return xproto.ChangePropertyChecked(X, xproto.PropModeReplace, win, propAtom, typ, format,
		uint32(len(data)/(int(format)/8)), data).Check()
}
