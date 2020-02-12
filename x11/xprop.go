package x11

import (
	"fmt"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func (xc *Connection) getProp(win xproto.Window, name string) (*xproto.GetPropertyReply, error) {
	atom := xc.Atom(name)
	cookie := xproto.GetProperty(xc.conn, false, win, atom, xproto.GetPropertyTypeAny, 0, 64)
	reply, err := cookie.Reply()
	if err != nil {
		return nil, fmt.Errorf("error retrieving property %q on window %d: %v", name, win, err)
	}
	if reply == nil || reply.Format == 0 {
		return nil, fmt.Errorf("no such property %q on window %d", name, win)
	}
	return reply, nil
}

func (xc *Connection) getProps32(win xproto.Window, name string) ([]uint32, error) {
	reply, err := xc.getProp(win, name)
	if err != nil {
		return nil, err
	}
	vals := make([]uint32, 0)
	for v := reply.Value; len(v) >= 4; v = v[4:] {
		vals = append(vals, uint32(v[0])|uint32(v[1])<<8|uint32(v[2])<<16|uint32(v[3])<<24)
	}
	return vals, nil
}

func (xc *Connection) changeProp32(win xproto.Window, prop string, typ xproto.Atom, data ...uint32) error {
	buf := make([]byte, len(data)*4)
	for i, datum := range data {
		xgb.Put32(buf[(i*4):], datum)
	}
	return xc.changeProp(win, 32, prop, typ, buf)
}

func (xc *Connection) changeProp(win xproto.Window, format byte, prop string, typ xproto.Atom, data []byte) error {
	propAtom := xc.Atom(prop)
	return xproto.ChangePropertyChecked(xc.conn, xproto.PropModeReplace, win, propAtom, typ, format,
		uint32(len(data)/(int(format)/8)), data).Check()
}
