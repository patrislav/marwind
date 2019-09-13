package keysym

import (
	"errors"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Keymap [256][]xproto.Keysym

const (
	loKey = 8
	hiKey = 255
)

func LoadKeyMapping(xc *xgb.Conn) (*Keymap, error) {
	m := xproto.GetKeyboardMapping(xc, loKey, hiKey-loKey+1)
	reply, err := m.Reply()
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, errors.New("could not load keysym map")
	}

	var keymap Keymap
	for i := 0; i < hiKey-loKey+1; i++ {
		keymap[loKey+i] = reply.Keysyms[i*int(reply.KeysymsPerKeycode) : (i+1)*int(reply.KeysymsPerKeycode)]
	}
	return &keymap, nil
}
