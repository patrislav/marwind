package x11

import (
	"image"

	"github.com/BurntSushi/xgbutil/xgraphics"
)

type Dimensions struct {
	Top, Left, Right, Bottom uint32
}

func (xc *Connection) NewImage(rect image.Rectangle) *xgraphics.Image {
	return xgraphics.New(xc.util, rect)
}
