package client

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"image"
	"testing"
)

type mockReparented struct {
	w, p xproto.Window
	x, y int16
}

type mockX11 struct {
	t *testing.T

	mappedWins     []xproto.Window
	unmappedWins   []xproto.Window
	destroyedWins  []xproto.Window
	reparentedWins []mockReparented
}

func (mx *mockX11) GetRootWindow() xproto.Window {
	return 0
}
func (mx *mockX11) CreateWindow(
	parent xproto.Window,
	x int16, y int16, width uint16, height uint16,
	borderWidth uint16,
	class uint16, valueMask uint32, valueList []uint32,
) (xproto.Window, error) {
	return 1, nil
}

func (mx *mockX11) MapWindow(window xproto.Window) error {
	return nil
}
func (mx *mockX11) UnmapWindow(window xproto.Window) error {
	return nil
}
func (mx *mockX11) DestroyWindow(window xproto.Window) error {
	return nil
}
func (mx *mockX11) ReparentWindow(window, parent xproto.Window, x, y int16) error {
	mx.reparentedWins = append(mx.reparentedWins, mockReparented{
		w: window,
		p: parent,
		x: x,
		y: y,
	})
	return nil
}

func (mx *mockX11) GetWindowTitle(window xproto.Window) (string, error) {
	return "", nil
}
func (mx *mockX11) Atom(name string) xproto.Atom {
	return 0
}

func (mx *mockX11) NewImage(rect image.Rectangle) *xgraphics.Image {
	return nil
}
