package wm

import (
	"github.com/BurntSushi/freetype-go/freetype"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/patrislav/marwind/x11"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	"image/draw"
)

type titlebarConfig struct {
	height      uint8
	borderWidth uint8
	bgColor     uint32
	fontColor   uint32
	fontSize    float64
}

type titlebar struct {
	frame  *frame
	win    xproto.Window
	config *titlebarConfig
}

func newTitlebar(f *frame, config *titlebarConfig) *titlebar {
	return &titlebar{
		config: config,
		frame:  f,
		win:    f.parent,
	}
}

func (t *titlebar) draw() error {
	title := t.frame.client.title
	width := t.frame.geom.W
	bg := color.RGBA{
		A: uint8((t.config.bgColor & 0xFF000000) >> 24),
		R: uint8((t.config.bgColor & 0x00FF0000) >> 16),
		G: uint8((t.config.bgColor & 0x0000FF00) >> 8),
		B: uint8(t.config.bgColor & 0x000000FF),
	}
	fg := color.RGBA{
		A: uint8((t.config.fontColor & 0xFF000000) >> 24),
		R: uint8((t.config.fontColor & 0x00FF0000) >> 16),
		G: uint8((t.config.fontColor & 0x0000FF00) >> 8),
		B: uint8(t.config.fontColor & 0x000000FF),
	}

	// title should never be zero-length
	if len(title) == 0 {
		title = " "
	}

	img := xgraphics.New(x11.XUtil, image.Rect(0, 0, int(width), int(t.config.height)))
	defer img.Destroy()
	img.ForExp(func(x, y int) (uint8, uint8, uint8, uint8) {
		return bg.R, bg.G, bg.B, bg.A
	})

	font, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		return err
	}

	// Over estimate the extents
	ew, eh := xgraphics.Extents(font, t.config.fontSize, title)

	// Create an image using the overestimated extents
	text := xgraphics.New(x11.XUtil, image.Rect(0, 0, ew, eh))
	defer text.Destroy()
	text.ForExp(func(x, y int) (uint8, uint8, uint8, uint8) {
		return bg.R, bg.G, bg.B, bg.A
	})
	_, _, err = text.Text(0, 0, fg, t.config.fontSize, font, title)
	if err != nil {
		return err
	}

	bounds := text.Bounds().Size()
	w, h := bounds.X, bounds.Y
	x := int(width/2) - w/2
	y := int(t.config.height/2) - h/2
	dstRect := image.Rect(x, y, x+w, y+h)
	draw.Draw(img, dstRect, text, image.ZP, draw.Src)

	if err := img.CreatePixmap(); err != nil {
		return err
	}
	img.XDraw()
	img.XExpPaint(t.win, int(t.config.borderWidth), int(t.config.borderWidth))
	return nil
}
