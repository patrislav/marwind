package wm

import (
	"fmt"
	"sort"

	"github.com/patrislav/marwind/x11"
)

type dockArea uint8

const (
	dockAreaTop    dockArea = 0
	dockAreaBottom dockArea = 1
)

type output struct {
	geom       x11.Geom
	workspaces []*workspace
	activeWs   *workspace
	dockAreas  [2][]*frame
}

// newOutput creates a new output from the given geometry
func newOutput(geom x11.Geom) *output {
	return &output{geom: geom}
}

// addWorkspace appends the workspace to this output, sorting them,
// and setting the activeWs if it's currently nil
func (o *output) addWorkspace(ws *workspace) error {
	ws.setOutput(o)
	o.workspaces = append(o.workspaces, ws)
	sort.Slice(o.workspaces, func(i, j int) bool {
		return o.workspaces[i].id < o.workspaces[j].id
	})
	if o.activeWs == nil {
		o.activeWs = ws
		return ws.show()
	}
	return nil
}

// addDock appends the frame as a dock of this output
func (o *output) addDock(f *frame) error {
	struts, err := x11.GetWindowStruts(f.client.window)
	if err != nil {
		return fmt.Errorf("failed to get struts: %v", err)
	}
	var area dockArea
	switch {
	case struts.Top > struts.Bottom:
		area = dockAreaTop
		f.height = struts.Top
	case struts.Bottom > struts.Top:
		area = dockAreaBottom
		f.height = struts.Bottom
	default:
		return fmt.Errorf("could not determine the dock position")
	}
	o.dockAreas[area] = append(o.dockAreas[area], f)
	// TODO map the dock
	o.updateTiling()
	return f.doMap()
}

// dockHeight returns the height of the entire dock area
func (o *output) dockHeight(area dockArea) uint32 {
	var height uint32
	for _, f := range o.dockAreas[area] {
		height += f.height
	}
	return height
}

func (o *output) workspaceArea() x11.Geom {
	top := o.dockHeight(dockAreaTop)
	bottom := o.dockHeight(dockAreaBottom)
	return x11.Geom{
		X: o.geom.X,
		Y: o.geom.Y + top,
		W: o.geom.W,
		H: o.geom.H - top - bottom,
	}
}

func (o *output) deleteFrame(frm *frame) bool {
	for area := range o.dockAreas {
		for i, f := range o.dockAreas[area] {
			if frm == f {
				o.dockAreas[area] = append(o.dockAreas[area][:i], o.dockAreas[area][i+1:]...)
				o.updateTiling()
				return true
			}
		}
	}
	for _, ws := range o.workspaces {
		if ws.deleteFrame(frm) {
			return true
		}
	}
	return false
}

func (o *output) updateTiling() {
	for _, ws := range o.workspaces {
		ws.updateTiling()
	}
}
