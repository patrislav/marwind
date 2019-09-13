package manager

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind-wm/container"
)

func (m *Manager) renderWorkspace(ws *container.Workspace) error {
	var err error
	var gap uint32
	if ws.HasGaps() {
		gap = m.config.InnerGap
	}
	startX := ws.Rect().X
	for _, col := range ws.Columns() {
		rect := container.Rect{
			X: startX,
			Y: ws.Rect().Y,
			W: col.Width(),
			H: ws.Rect().H,
		}
		err = m.renderColumn(col, rect, gap)
		startX += col.Width()
	}
	return err
}

func (m *Manager) renderColumn(col *container.Column, rect container.Rect, gap uint32) error {
	var err error
	startY := rect.Y
	for _, frame := range col.Frames() {
		rect := container.Rect{
			X: rect.X,
			Y: startY,
			W: rect.W,
			H: frame.Height(),
		}
		err = m.renderFrame(frame, rect, gap)
		startY += frame.Height()
	}
	return err
}

func (m *Manager) renderFrame(frame *container.Frame, rect container.Rect, gap uint32) error {
	mask := uint16(xproto.ConfigWindowX | xproto.ConfigWindowY | xproto.ConfigWindowWidth | xproto.ConfigWindowHeight)
	values := []uint32{
		rect.X + gap,
		rect.Y + gap,
		rect.W - gap*2,
		rect.H - gap*2,
	}
	return xproto.ConfigureWindowChecked(m.xc, frame.Window(), mask, values).Check()
}
