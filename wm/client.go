package wm

import (
	"github.com/BurntSushi/xgb/xproto"
)

// client wraps a single X client (window)
type client struct {
	frame  *frame
	window xproto.Window
}
