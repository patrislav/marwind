package x11

import (
	"github.com/BurntSushi/xgbutil/xgraphics"
	"image"
)

type Dimensions struct {
	Top, Left, Right, Bottom uint32
}

func (xc *Connection) NewImage(rect image.Rectangle) *xgraphics.Image {
	return xgraphics.New(xc.util, rect)
}
