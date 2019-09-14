package container

import (
	"fmt"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type MoveDirection uint8

const (
	MoveLeft MoveDirection = iota
	MoveRight
	MoveUp
	MoveDown
)

type WorkspaceConfig struct {
	Gap uint32
}

type Workspace struct {
	output  *Output
	columns []*Column
	config  WorkspaceConfig
}

func NewWorkspace(config WorkspaceConfig) *Workspace {
	return &Workspace{config: config}
}

func (ws *Workspace) AddWindow(xc *xgb.Conn, win xproto.Window) error {
	frame, err := ManageWindow(xc, win)
	if err != nil {
		return err
	}
	var col *Column
	if len(ws.columns) < 2 {
		col = ws.createColumn(false)
	}
	if col == nil {
		col = ws.columns[len(ws.columns)-1]
	}
	col.AddFrame(frame, nil)
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
	return ws.output.rect
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
