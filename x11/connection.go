package x11

import (
	"github.com/BurntSushi/xgb"
)

var X *xgb.Conn

func InitConnection() error {
	var err error
	X, err = xgb.NewConn()
	if err != nil {
		return err
	}
	return nil
}
