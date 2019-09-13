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
		col = ws.createColumn()
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
	// TODO: implement MoveWindow
	return nil
}

func (ws *Workspace) Columns() []*Column { return ws.columns }

func (ws *Workspace) Rect() Rect {
	r := ws.output.rect
	return Rect{
		X: r.X + ws.config.Gap,
		Y: r.Y + ws.config.Gap,
		W: r.W - ws.config.Gap*2,
		H: r.H - ws.config.Gap*2,
	}
}

func (ws *Workspace) setOutput(output *Output) {
	ws.output = output
}

func (ws *Workspace) createColumn() *Column {
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
	ws.columns = append(ws.columns, col)
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
