package container

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind/x11"
)

type DockArea uint8

const (
	DockAreaTop    DockArea = 0
	DockAreaBottom DockArea = 1
)

type Output struct {
	rect       Rect
	workspaces []*Workspace
	currentWs  *Workspace
	dockAreas  [2][]*Frame
}

func NewOutput(rect Rect) *Output {
	return &Output{rect: rect}
}

func (o *Output) Rect() Rect                        { return o.rect }
func (o *Output) DockFrames(area DockArea) []*Frame { return o.dockAreas[area] }
func (o *Output) CurrentWorkspace() *Workspace      { return o.currentWs }
func (o *Output) Workspaces() []*Workspace          { return o.workspaces }

func (o *Output) SwitchWorkspace(next *Workspace) error {
	if next == o.currentWs {
		return nil
	}
	check := o.FindWorkspace(func(ws *Workspace) bool { return ws == next })
	if check == nil {
		return fmt.Errorf("workspace is not a part of this output")
	}
	err := o.currentWs.hide()
	if err != nil {
		return err
	}
	err = next.show()
	if err != nil {
		return err
	}
	o.currentWs = next
	return nil
}

func (o *Output) AddWorkspace(ws *Workspace) error {
	ws.setOutput(o)
	o.workspaces = append(o.workspaces, ws)
	if o.currentWs == nil {
		o.currentWs = ws
		return ws.show()
	}
	return nil
}

func (o *Output) AddDock(f *Frame) error {
	struts, err := x11.GetWindowStruts(f.window)
	if err != nil {
		return err
	}
	var area DockArea
	switch {
	case struts.Top > struts.Bottom:
		area = DockAreaTop
		f.height = struts.Top
	case struts.Bottom > struts.Top:
		area = DockAreaBottom
		f.height = struts.Bottom
	default:
		return fmt.Errorf("could not determine the dock position")
	}
	o.dockAreas[area] = append(o.dockAreas[area], f)
	o.updateTiling()
	return f.Map()
}

func (o *Output) DeleteWindow(win xproto.Window) bool {
	for area := range o.dockAreas {
		for i, f := range o.dockAreas[area] {
			if f.window == win {
				o.dockAreas[area] = append(o.dockAreas[area][:i], o.dockAreas[area][i+1:]...)
				o.updateTiling()
				return true
			}
		}
	}
	for _, ws := range o.workspaces {
		if err := ws.DeleteWindow(win); err == nil {
			return true
		}
	}
	return false
}

func (o *Output) updateTiling() {
	for _, ws := range o.workspaces {
		ws.UpdateTiling()
	}
}

func (o *Output) workspaceRect() Rect {
	top := o.DockHeight(DockAreaTop)
	bottom := o.DockHeight(DockAreaBottom)
	return Rect{
		X: o.rect.X,
		Y: o.rect.Y + top,
		W: o.rect.W,
		H: o.rect.H - top - bottom,
	}
}

func (o *Output) DockHeight(area DockArea) uint32 {
	var height uint32
	for _, f := range o.dockAreas[area] {
		height += f.height
	}
	return height
}

func (o *Output) FindFrame(predicate func(*Frame) bool) *Frame {
	for _, ws := range o.workspaces {
		f := ws.findFrame(predicate)
		if f != nil {
			return f
		}
	}
	return nil
}

func (o *Output) FindWorkspace(predicate func(*Workspace) bool) *Workspace {
	for _, ws := range o.workspaces {
		if predicate(ws) {
			return ws
		}
	}
	return nil
}
