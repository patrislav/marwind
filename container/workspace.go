package container

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"
)

type MoveDirection uint8

const (
	MoveLeft MoveDirection = iota
	MoveRight
	MoveUp
	MoveDown
)

type ResizeDirection uint8

const (
	ResizeVert ResizeDirection = iota
	ResizeHoriz
)

type WorkspaceConfig struct {
	Gap uint32
}

type Workspace struct {
	ID      uint8
	Name    string
	output  *Output
	columns []*Column
	config  WorkspaceConfig
}

func NewWorkspace(id uint8, config WorkspaceConfig) *Workspace {
	return &Workspace{config: config, ID: id, Name: fmt.Sprintf("%d", id+1)}
}

func (ws *Workspace) AddFrame(frame *Frame) error {
	var col *Column
	if len(ws.columns) < 2 {
		col = ws.createColumn(false)
	}
	if col == nil {
		col = ws.columns[len(ws.columns)-1]
	}
	col.AddFrame(frame, nil)
	if ws.output.currentWs == ws {
		return frame.Map()
	}
	return nil
}

func (ws *Workspace) DeleteWindow(win xproto.Window) error {
	frame := ws.findFrame(func(f *Frame) bool { return f.window == win })
	if frame == nil {
		return fmt.Errorf("could not find frame with window %d", win)
	}
	col := frame.col
	col.DeleteFrame(frame)
	if len(col.frames) == 0 {
		ws.deleteColumn(col)
	}
	return nil
}

func (ws *Workspace) MoveWindow(win xproto.Window, dir MoveDirection) error {
	frame := ws.findFrame(func(f *Frame) bool { return f.window == win })
	if frame == nil {
		return nil
	}
	switch dir {
	case MoveLeft:
		i := ws.findColumnIndex(func(c *Column) bool { return c == frame.col })
		origCol := frame.col
		origCol.DeleteFrame(frame)
		if i == 0 {
			col := ws.createColumn(true)
			col.AddFrame(frame, nil)
		} else {
			col := ws.columns[i-1]
			col.AddFrame(frame, nil)
		}
		if len(origCol.frames) == 0 {
			ws.deleteColumn(origCol)
		}
	case MoveRight:
		i := ws.findColumnIndex(func(c *Column) bool { return c == frame.col })
		origCol := frame.col
		origCol.DeleteFrame(frame)
		if i == len(ws.columns)-1 {
			col := ws.createColumn(false)
			col.AddFrame(frame, nil)
		} else {
			col := ws.columns[i+1]
			col.AddFrame(frame, nil)
		}
		if len(origCol.frames) == 0 {
			ws.deleteColumn(origCol)
		}
	case MoveUp:
		col := frame.col
		i := col.findFrameIndex(func(f *Frame) bool { return f == frame })
		if i > 0 {
			other := col.frames[i-1]
			col.frames[i-1] = frame
			col.frames[i] = other
		}
	case MoveDown:
		col := frame.col
		i := col.findFrameIndex(func(f *Frame) bool { return f == frame })
		if i < len(col.frames)-1 {
			other := col.frames[i+1]
			col.frames[i+1] = frame
			col.frames[i] = other
		}
	}
	return nil
}

func (ws *Workspace) ResizeWindow(win xproto.Window, dir ResizeDirection, pct int) error {
	frame := ws.findFrame(func(f *Frame) bool { return f.window == win })
	if frame == nil {
		return nil
	}
	switch dir {
	case ResizeHoriz:
		if len(ws.columns) < 2 {
			return nil
		}
		min := uint32(float32(ws.Rect().W) * 0.1)
		dwFull := int(float32(ws.Rect().W) * (float32(pct) / 100))
		if uint32(int(frame.col.width)+dwFull) < min {
			return nil
		}
		dwPart := dwFull/len(ws.columns) - 1
		dwFinal := 0
		for _, col := range ws.columns {
			if col != frame.col {
				next := uint32(int(col.width) - dwPart)
				if next >= min {
					col.width = next
					dwFinal += dwPart
				}
			}
		}
		frame.col.width = uint32(int(frame.col.width) + dwFinal)
	case ResizeVert:
		col := frame.col
		if len(col.frames) < 2 {
			return nil
		}
		min := uint32(float32(ws.Rect().H) * 0.1)
		dhFull := int(float32(ws.Rect().H) * (float32(pct) / 100))
		if uint32(int(frame.height)+dhFull) < min {
			return nil
		}
		dhPart := dhFull/len(col.frames) - 1
		dhFinal := 0
		for _, f := range col.frames {
			if f != frame {
				next := uint32(int(f.height) - dhPart)
				if next >= min {
					f.height = next
					dhFinal += dhPart
				}
			}
		}
		frame.height = uint32(int(frame.height) + dhFinal)
	}
	return nil
}

func (ws *Workspace) Columns() []*Column { return ws.columns }

func (ws *Workspace) Rect() Rect {
	r := ws.FullRect()
	return Rect{
		X: r.X + ws.config.Gap,
		Y: r.Y + ws.config.Gap,
		W: r.W - ws.config.Gap*2,
		H: r.H - ws.config.Gap*2,
	}
}

func (ws *Workspace) FullRect() Rect {
	return ws.output.workspaceRect()
}

func (ws *Workspace) HasWindow(win xproto.Window) bool {
	frame := ws.findFrame(func(f *Frame) bool { return f.window == win })
	return frame != nil
}

// GetOnlyFrame returns a pointer to the frame if there's only one in the workspace and nil otherwise
func (ws *Workspace) GetOnlyFrame() *Frame {
	if ws.countAllFrames() == 1 {
		return ws.columns[0].frames[0]
	}
	return nil
}

func (ws *Workspace) UpdateTiling() {
	for _, col := range ws.columns {
		col.UpdateTiling()
	}
}

func (ws *Workspace) IsVisible() bool { return ws.output.currentWs == ws }
func (ws *Workspace) Output() *Output { return ws.output }

func (ws *Workspace) setOutput(output *Output) {
	ws.output = output
}

func (ws *Workspace) createColumn(start bool) *Column {
	wsWidth := ws.Rect().W
	origLen := len(ws.columns)
	col := &Column{ws: ws, width: ws.Rect().W / uint32(origLen+1)}
	if origLen > 0 {
		col.width = wsWidth / uint32(origLen+1)
		remWidth := float32(wsWidth - col.width)
		leftWidth := uint32(remWidth)
		for _, c := range ws.columns {
			c.width = uint32((float32(c.width) / float32(wsWidth)) * remWidth)
			leftWidth -= c.width
		}
		if leftWidth != 0 {
			col.width += leftWidth
		}
	} else {
		col.width = wsWidth
	}
	if start {
		ws.columns = append([]*Column{col}, ws.columns...)
	} else {
		ws.columns = append(ws.columns, col)
	}
	return col
}

func (ws *Workspace) deleteColumn(col *Column) {
	i := ws.findColumnIndex(func(c *Column) bool { return c == col })
	if i < 0 {
		return
	}
	wsWidth := ws.Rect().W
	// TODO: assign the widths proportional to the original width/totalWidth ratio
	// origLen = len(ws.columns)
	ws.columns = append(ws.columns[:i], ws.columns[i+1:]...)
	for _, c := range ws.columns {
		c.width = wsWidth / uint32(len(ws.columns))
	}
}

func (ws *Workspace) findFrame(predicate func(*Frame) bool) *Frame {
	for _, col := range ws.columns {
		idx := col.findFrameIndex(predicate)
		if idx >= 0 {
			return col.frames[idx]
		}
	}
	return nil
}

func (ws *Workspace) findColumnIndex(predicate func(*Column) bool) int {
	for i, col := range ws.columns {
		if predicate(col) {
			return i
		}
	}
	return -1
}

func (ws *Workspace) countAllFrames() int {
	count := 0
	for _, col := range ws.columns {
		count += len(col.frames)
	}
	return count
}

func (ws *Workspace) hide() error {
	var err error
	for _, col := range ws.columns {
		for _, f := range col.frames {
			err = f.Unmap()
		}
	}
	return err
}

func (ws *Workspace) show() error {
	var err error
	for _, col := range ws.columns {
		for _, f := range col.frames {
			err = f.Map()
		}
	}
	return err
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
