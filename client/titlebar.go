package client

import (
	"github.com/BurntSushi/freetype-go/freetype"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	"image/draw"
)

func (c *Client) drawTitlebar() error {
	width := c.geom.W
	bg := color.RGBA{
		A: uint8((c.cfg.BgColor & 0xFF000000) >> 24),
		R: uint8((c.cfg.BgColor & 0x00FF0000) >> 16),
		G: uint8((c.cfg.BgColor & 0x0000FF00) >> 8),
		B: uint8(c.cfg.BgColor & 0x000000FF),
	}
	fg := color.RGBA{
		A: uint8((c.cfg.FontColor & 0xFF000000) >> 24),
		R: uint8((c.cfg.FontColor & 0x00FF0000) >> 16),
		G: uint8((c.cfg.FontColor & 0x0000FF00) >> 8),
		B: uint8(c.cfg.FontColor & 0x000000FF),
	}

	// title should never be zero-length
	if len(c.title) == 0 {
		c.title = " "
	}

	img := c.x11.NewImage(image.Rect(0, 0, int(width), int(c.cfg.TitlebarHeight)))
	defer img.Destroy()
	img.ForExp(func(x, y int) (uint8, uint8, uint8, uint8) {
		return bg.R, bg.G, bg.B, bg.A
	})

	font, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		return err
	}

	// Over estimate the extents
	ew, eh := xgraphics.Extents(font, c.cfg.FontSize, c.title)

	// Create an image using the overestimated extents
	text := c.x11.NewImage(image.Rect(0, 0, ew, eh))
	defer text.Destroy()
	text.ForExp(func(x, y int) (uint8, uint8, uint8, uint8) {
		return bg.R, bg.G, bg.B, bg.A
	})
	_, _, err = text.Text(0, 0, fg, c.cfg.FontSize, font, c.title)
	if err != nil {
		return err
	}

	bounds := text.Bounds().Size()
	w, h := bounds.X, bounds.Y
	x := int(width/2) - w/2
	y := int(c.cfg.TitlebarHeight/2) - h/2
	dstRect := image.Rect(x, y, x+w, y+h)
	draw.Draw(img, dstRect, text, image.ZP, draw.Src)

	if err := img.CreatePixmap(); err != nil {
		return err
	}
	img.XDraw()
	img.XExpPaint(c.parent, int(c.cfg.BorderWidth), int(c.cfg.BorderWidth))
	return nil
}
