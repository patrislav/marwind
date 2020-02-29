package x11

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
)

type Connection struct {
	conn   *xgb.Conn
	util   *xgbutil.XUtil
	screen xproto.ScreenInfo
	atoms  map[string]xproto.Atom
}

func Connect() (*Connection, error) {
	atoms := make(map[string]xproto.Atom)
	xconn, err := xgb.NewConn()
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	xutil, err := xgbutil.NewConnXgb(xconn)
	if err != nil {
		return nil, fmt.Errorf("failed to create XUtil connection: %w", err)
	}
	return &Connection{conn: xconn, util: xutil, atoms: atoms}, nil
}

func (xc *Connection) X() *xgb.Conn              { return xc.conn }
func (xc *Connection) Screen() xproto.ScreenInfo { return xc.screen }

func (xc *Connection) Init() error {
	conninfo := xproto.Setup(xc.conn)
	if conninfo == nil {
		return errors.New("could not parse X connection info")
	}
	if len(conninfo.Roots) != 1 {
		return errors.New("wrong number of roots, possibly xinerama did not initialize properly")
	}
	xc.screen = conninfo.Roots[0]

	err := xc.setHints()
	if err != nil {
		return err
	}
	err = xc.initDesktop()
	if err != nil {
		return err
	}
	return nil
}

func (xc *Connection) Close() {
	xc.conn.Close()
}
