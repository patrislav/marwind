package x11

import (
	"errors"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
)

var (
	X      *xgb.Conn
	Screen xproto.ScreenInfo
)

func CreateConnection() error {
	var err error
	X, err = xgb.NewConn()
	if err != nil {
		return err
	}
	return nil
}

func InitConnection() error {
	if err := xinerama.Init(X); err != nil {
		return err
	}

	conninfo := xproto.Setup(X)
	if conninfo == nil {
		return errors.New("could not parse X connection info")
	}
	if len(conninfo.Roots) != 1 {
		return errors.New("wrong number of roots, did xinerama initialize properly?")
	}
	Screen = conninfo.Roots[0]

	err := setHints()
	if err != nil {
		return err
	}
	err = initDesktop()
	if err != nil {
		return err
	}
	return nil
}
