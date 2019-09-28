package wm

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind/x11"
)

func (wm *WM) renderOutput(o *output) error {
	var err error
	if e := wm.renderDock(o, dockAreaTop); e != nil {
		err = e
	}
	if e := wm.renderDock(o, dockAreaBottom); e != nil {
		err = e
	}
	if e := wm.renderWorkspace(o.activeWs); e != nil {
		err = e
	}
	return err
}

func (wm *WM) renderDock(o *output, area dockArea) error {
	var err error
	var y uint32
	switch area {
	case dockAreaTop:
		y = o.geom.Y
	case dockAreaBottom:
		y = o.geom.H - o.dockHeight(area)
	}
	for _, f := range o.dockAreas[area] {
		geom := x11.Geom{
			X: o.geom.X,
			Y: y,
			W: o.geom.W,
			H: f.height,
		}
		err = wm.renderFrame(f, geom)
		y += geom.H
	}
	return err
}

func (wm *WM) renderWorkspace(ws *workspace) error {
	var err error
	area := ws.output.workspaceArea()
	x := area.X
	for _, col := range ws.columns {
		geom := x11.Geom{
			X: x,
			Y: area.Y,
			W: col.width,
			H: area.H,
		}
		if e := wm.renderColumn(col, geom); e != nil {
			err = e
		}
		x += col.width
	}
	return err
}

func (wm *WM) renderColumn(col *column, geom x11.Geom) error {
	var err error
	y := geom.Y
	for _, f := range col.frames {
		fg := x11.Geom{
			X: geom.X,
			Y: y,
			W: geom.W,
			H: f.height,
		}
		if e := wm.renderFrame(f, fg); e != nil {
			err = e
		}
		y += f.height
	}
	return err
}

func (wm *WM) renderFrame(f *frame, geom x11.Geom) error {
	if !f.mapped {
		return nil
	}
	f.geom = geom
	mask := uint16(xproto.ConfigWindowX | xproto.ConfigWindowY | xproto.ConfigWindowWidth | xproto.ConfigWindowHeight)
	values := []uint32{
		geom.X,
		geom.Y,
		geom.W,
		geom.H,
	}
	err := xproto.ConfigureWindowChecked(x11.X, f.parent, mask, values).Check()
	if err != nil {
		return err
	}
	return xproto.ConfigureWindowChecked(x11.X, f.client.window, mask, []uint32{0, 0, geom.W, geom.H}).Check()
}
