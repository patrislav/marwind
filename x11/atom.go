package x11

import (
	"github.com/BurntSushi/xgb/xproto"
	"log"
)

// Atom returns the X11 atom of the given name
func (xc *Connection) Atom(name string) xproto.Atom {
	if atom, ok := xc.atoms[name]; ok {
		return atom
	}
	reply, err := xproto.InternAtom(xc.conn, false, uint16(len(name)), name).Reply()
	if err != nil {
		log.Fatal(err)
	}
	if reply == nil {
		return 0
	}
	xc.atoms[name] = reply.Atom
	return reply.Atom
}
