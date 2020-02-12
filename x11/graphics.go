package x11

import (
	"github.com/BurntSushi/xgbutil/xgraphics"
	"image"
)

func (xc *Connection) NewImage(rect image.Rectangle) *xgraphics.Image {
	return xgraphics.New(xc.util, rect)
}
