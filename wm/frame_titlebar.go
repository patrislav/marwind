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

type titlebar struct {
	frame *frame
	win   xproto.Window
	bg    color.RGBA
}

func newTitlebar(f *frame, bg uint32) *titlebar {
	clr := color.RGBA{
		A: uint8((bg & 0xFF000000) >> 24),
		R: uint8((bg & 0x00FF0000) >> 16),
		G: uint8((bg & 0x0000FF00) >> 8),
		B: uint8(bg & 0x000000FF),
	}
	return &titlebar{frame: f, win: f.parent, bg: clr}
}

func (t *titlebar) draw() error {
	title := t.frame.client.title
	width := t.frame.geom.W
	height := 18

	// title should never be zero-length
	if len(title) == 0 {
		title = " "
	}

	img := xgraphics.New(x11.XUtil, image.Rect(0, 0, int(width), int(height)))
	defer img.Destroy()
	img.ForExp(func(x, y int) (uint8, uint8, uint8, uint8) {
		return t.bg.R, t.bg.G, t.bg.B, t.bg.A
	})

	fgCol := color.Black
	size := float64(12)
	font, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		return err
	}

	// Over estimate the extents
	ew, eh := xgraphics.Extents(font, size, title)

	// Create an image using the overestimated extents
	text := xgraphics.New(x11.XUtil, image.Rect(0, 0, ew, eh))
	defer text.Destroy()
	text.ForExp(func(x, y int) (uint8, uint8, uint8, uint8) {
		return t.bg.R, t.bg.G, t.bg.B, t.bg.A
	})
	_, _, err = text.Text(0, 0, fgCol, size, font, title)
	if err != nil {
		return err
	}

	bounds := text.Bounds().Size()
	w, h := bounds.X, bounds.Y
	x := int(width/2) - w/2
	y := height/2 - h/2
	dstRect := image.Rect(x, y, x+w, y+h)
	draw.Draw(img, dstRect, text, image.ZP, draw.Src)

	if err := img.CreatePixmap(); err != nil {
		return err
	}
	img.XDraw()
	img.XExpPaint(t.win, 1, 1)
	return nil
}
