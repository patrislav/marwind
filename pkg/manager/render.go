package manager

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind/pkg/container"
	"github.com/patrislav/marwind/pkg/x11"
)

func (m *Manager) renderOutput(o *container.Output) error {
	var err error
	if e := m.renderDock(o, container.DockAreaTop); e != nil {
		err = e
	}
	if e := m.renderDock(o, container.DockAreaBottom); e != nil {
		err = e
	}
	if e := m.renderWorkspace(o.CurrentWorkspace()); e != nil {
		err = e
	}
	return err
}

func (m *Manager) renderDock(o *container.Output, area container.DockArea) error {
	var err error
	var y uint32
	switch area {
	case container.DockAreaTop:
		y = o.Rect().Y
	case container.DockAreaBottom:
		y = o.Rect().H - o.DockHeight(area)
	}
	for _, f := range o.DockFrames(area) {
		rect := container.Rect{
			X: o.Rect().X,
			Y: y,
			W: o.Rect().W,
			H: f.Height(),
		}
		err = m.renderFrame(f, rect)
		y += rect.H
	}
	return err
}

func (m *Manager) renderWorkspace(ws *container.Workspace) error {
	var err error
	onlyFrame := ws.GetOnlyFrame()
	if onlyFrame != nil {
		return m.renderFrame(onlyFrame, ws.FullRect())
	}
	startX := ws.Rect().X
	for _, col := range ws.Columns() {
		rect := container.Rect{
			X: startX,
			Y: ws.Rect().Y,
			W: col.Width(),
			H: ws.Rect().H,
		}
		err = m.renderColumn(col, rect, m.config.InnerGap)
		startX += col.Width()
	}
	return err
}

func (m *Manager) renderColumn(col *container.Column, rect container.Rect, gap uint32) error {
	var err error
	startY := rect.Y
	for _, frame := range col.Frames() {
		rect := container.Rect{
			X: rect.X + gap,
			Y: startY + gap,
			W: rect.W - gap*2,
			H: frame.Height() - gap*2,
		}
		err = m.renderFrame(frame, rect)
		startY += frame.Height()
	}
	return err
}

func (m *Manager) renderFrame(frame *container.Frame, rect container.Rect) error {
	frame.Rect = rect
	mask := uint16(xproto.ConfigWindowX | xproto.ConfigWindowY | xproto.ConfigWindowWidth | xproto.ConfigWindowHeight)
	values := []uint32{
		rect.X,
		rect.Y,
		rect.W,
		rect.H,
	}
	return xproto.ConfigureWindowChecked(x11.X, frame.Window(), mask, values).Check()
}
