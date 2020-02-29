package wm

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind/client"
	"github.com/patrislav/marwind/x11"
)

type frame struct {
	col    *column
	cli    *client.Client
	height uint16
}

func (wm *WM) createFrame(win xproto.Window, typ client.Type) (*frame, error) {
	c, err := client.New(wm.xc, wm.windowConfig, win, typ)
	if err != nil {
		return nil, err
	}
	f := &frame{cli: c}

	return f, nil
}

func (f *frame) workspace() *workspace {
	if f.col != nil {
		return f.col.ws
	}
	return nil
}

func (wm *WM) getFrameDecorations(f *frame) x11.Dimensions {
	if f.cli.Parent() == 0 {
		return x11.Dimensions{Top: 0, Left: 0, Right: 0, Bottom: 0}
	}
	var bar uint32
	border := uint32(wm.config.BorderWidth)
	if wm.config.TitleBarHeight > 0 {
		bar = uint32(wm.config.TitleBarHeight) + 1
	}
	return x11.Dimensions{
		Top:    border + bar,
		Right:  border,
		Bottom: border,
		Left:   border,
	}
}
