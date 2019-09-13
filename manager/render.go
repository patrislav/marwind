package manager

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind-wm/container"
)

func (m *Manager) renderWorkspace(ws *container.Workspace) error {
	fmt.Println("rendering workspace")
	var err error
	startX := ws.Rect().X
	for _, col := range ws.Columns() {
		rect := container.Rect{
			X: startX,
			Y: ws.Rect().Y,
			W: col.Width(),
			H: ws.Rect().H,
		}
		err = m.renderColumn(col, rect)
		startX += col.Width()
	}
	return err
}

func (m *Manager) renderColumn(col *container.Column, rect container.Rect) error {
	fmt.Println("rendering column", rect)
	var err error
	startY := rect.Y
	for _, frame := range col.Frames() {
		rect := container.Rect{
			X: rect.X,
			Y: startY,
			W: rect.W,
			H: frame.Height(),
		}
		err = m.renderFrame(frame, rect)
		startY += frame.Height()
	}
	return err
}

func (m *Manager) renderFrame(frame *container.Frame, rect container.Rect) error {
	fmt.Println("rendering frame", rect)
	mask := uint16(xproto.ConfigWindowX | xproto.ConfigWindowY | xproto.ConfigWindowWidth | xproto.ConfigWindowHeight)
	values := []uint32{
		rect.X + m.config.InnerGap,
		rect.Y + m.config.InnerGap,
		rect.W - m.config.InnerGap*2,
		rect.H - m.config.InnerGap*2,
	}
	return xproto.ConfigureWindowChecked(m.xc, frame.Window(), mask, values).Check()
}
