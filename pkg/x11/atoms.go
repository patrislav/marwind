package x11

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"
)

var atoms = make(map[string]xproto.Atom)

func Atom(name string) xproto.Atom {
	if atom, ok := atoms[name]; ok {
		return atom
	}
	reply, err := xproto.InternAtom(X, false, uint16(len(name)), name).Reply()
	if err != nil {
		log.Fatal(err)
	}
	if reply == nil {
		return 0
	}
	atoms[name] = reply.Atom
	return reply.Atom
}
