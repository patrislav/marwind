package wm

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind/client"
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
	var y int16
	switch area {
	case dockAreaTop:
		y = o.geom.Y
	case dockAreaBottom:
		y = int16(o.geom.H - o.dockHeight(area))
	}
	for _, f := range o.dockAreas[area] {
		geom := client.Geom{
			X: o.geom.X,
			Y: y,
			W: o.geom.W,
			H: f.height,
		}
		err = wm.renderFrame(f, geom)
		y += int16(geom.H)
	}
	return err
}

func (wm *WM) renderWorkspace(ws *workspace) error {
	var err error
	if f := ws.singleFrame(); f != nil {
		return wm.renderFrame(f, ws.fullArea())
	}
	a := ws.area()
	x := a.X
	for _, col := range ws.columns {
		geom := client.Geom{
			X: x,
			Y: a.Y,
			W: col.width,
			H: a.H,
		}
		if e := wm.renderColumn(col, geom); e != nil {
			err = e
		}
		x += int16(col.width)
	}
	return err
}

func (wm *WM) renderColumn(col *column, geom client.Geom) error {
	var err error
	y := geom.Y
	gap := wm.config.InnerGap
	for _, f := range col.frames {
		fg := client.Geom{
			X: geom.X + int16(gap),
			Y: y + int16(gap),
			W: geom.W - gap*2,
			H: f.height - gap*2,
		}
		if e := wm.renderFrame(f, fg); e != nil {
			err = e
		}
		y += int16(f.height)
	}
	return err
}

func (wm *WM) renderFrame(f *frame, geom client.Geom) error {
	if !f.cli.Mapped() {
		return nil
	}
	f.cli.SetGeom(geom)
	mask := uint16(xproto.ConfigWindowX | xproto.ConfigWindowY | xproto.ConfigWindowWidth | xproto.ConfigWindowHeight)
	parentVals := []uint32{uint32(geom.X), uint32(geom.Y), uint32(geom.W), uint32(geom.H)}
	clientVals := parentVals
	if f.cli.Parent() != 0 {
		if err := xproto.ConfigureWindowChecked(wm.xc.X(), f.cli.Parent(), mask, parentVals).Check(); err != nil {
			return err
		}
		d := wm.getFrameDecorations(f)
		clientVals = []uint32{d.Left, d.Top, uint32(geom.W) - d.Left - d.Right, uint32(geom.H) - d.Top - d.Bottom}
	}
	if err := xproto.ConfigureWindowChecked(wm.xc.X(), f.cli.Window(), mask, clientVals).Check(); err != nil {
		return err
	}
	if err := wm.configureNotify(f); err != nil {
		return err
	}
	return nil
}

func (wm *WM) configureNotify(f *frame) error {
	// Hack for Java applications as described here:
	// https://stackoverflow.com/questions/31646544/xlib-reparenting-a-java-window-with-popups-properly-translated
	// TODO: when window decorations are added, this should change to include them
	geom := f.cli.Geom()
	if f.cli.Parent() != 0 {
		d := wm.getFrameDecorations(f)
		geom = client.Geom{
			X: geom.X + int16(d.Left),
			Y: geom.Y + int16(d.Top),
			W: geom.W - uint16(d.Left-d.Right),
			H: geom.H - uint16(d.Top-d.Bottom),
		}
	}
	ev := xproto.ConfigureNotifyEvent{
		Event:            f.cli.Window(),
		Window:           f.cli.Window(),
		X:                geom.X,
		Y:                geom.Y,
		Width:            geom.W,
		Height:           geom.H,
		BorderWidth:      0,
		AboveSibling:     0,
		OverrideRedirect: true,
	}
	evCookie := xproto.SendEventChecked(wm.xc.X(), false, f.cli.Window(), xproto.EventMaskStructureNotify, string(ev.Bytes()))
	if err := evCookie.Check(); err != nil {
		return err
	}
	return nil
}
