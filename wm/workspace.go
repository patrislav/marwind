package wm

import (
	"github.com/patrislav/marwind/x11"
)

type workspaceConfig struct {
	gap uint32
}

type workspace struct {
	id      uint8
	columns []*column
	output  *output
	config  workspaceConfig
}

func newWorkspace(id uint8, config workspaceConfig) *workspace {
	return &workspace{id: id, config: config}
}

func (ws *workspace) setOutput(o *output) {
	ws.output = o
}

// addFrame appends the given frame to the last column in the workspace
func (ws *workspace) addFrame(f *frame) error {
	var col *column
	if len(ws.columns) < 2 {
		col = ws.createColumn(false)
	}
	if col == nil {
		col = ws.columns[len(ws.columns)-1]
	}
	col.addFrame(f, nil)
	if ws.output.activeWs == ws {
		return f.doMap()
	}
	return nil
}

// deleteFrame deletes the frame from any column that contains it
func (ws *workspace) deleteFrame(f *frame) bool {
	if f.col == nil || f.col.ws != ws {
		return false
	}
	col := f.col
	col.deleteFrame(f)
	if len(col.frames) == 0 {
		ws.deleteColumn(col)
	}
	return true
}

// moveFrame changes the position of a frame within a column or moves it between columns
func (ws *workspace) moveFrame(f *frame, dir MoveDirection) error {
	switch dir {
	case MoveLeft:
		i := ws.findColumnIndex(func(c *column) bool { return c == f.col })
		origCol := f.col
		origCol.deleteFrame(f)
		if i == 0 {
			col := ws.createColumn(true)
			col.addFrame(f, nil)
		} else {
			col := ws.columns[i-1]
			col.addFrame(f, nil)
		}
		if len(origCol.frames) == 0 {
			ws.deleteColumn(origCol)
		}
	case MoveRight:
		i := ws.findColumnIndex(func(c *column) bool { return c == f.col })
		origCol := f.col
		origCol.deleteFrame(f)
		if i == len(ws.columns)-1 {
			col := ws.createColumn(false)
			col.addFrame(f, nil)
		} else {
			col := ws.columns[i+1]
			col.addFrame(f, nil)
		}
		if len(origCol.frames) == 0 {
			ws.deleteColumn(origCol)
		}
	case MoveUp:
		col := f.col
		i := col.findFrameIndex(func(frm *frame) bool { return f == frm })
		if i > 0 {
			other := col.frames[i-1]
			col.frames[i-1] = f
			col.frames[i] = other
		}
	case MoveDown:
		col := f.col
		i := col.findFrameIndex(func(frm *frame) bool { return f == frm })
		if i < len(col.frames)-1 {
			other := col.frames[i+1]
			col.frames[i+1] = f
			col.frames[i] = other
		}
	}
	return nil
}

// show maps all the frames of the workspace
func (ws *workspace) show() error {
	var err error
	for _, col := range ws.columns {
		for _, f := range col.frames {
			if e := f.doMap(); e != nil {
				err = e
			}
		}
	}
	return err
}

// hide unmaps all the frames of the workspace
func (ws *workspace) hide() error {
	var err error
	for _, col := range ws.columns {
		for _, f := range col.frames {
			if e := f.doUnmap(); e != nil {
				err = e
			}
		}
	}
	return err
}

// createColumn creates a new empty column either at the start (if the start argument is true)
// or the end of the workspace area.
func (ws *workspace) createColumn(start bool) *column {
	wsArea := ws.output.workspaceArea()
	wsWidth := wsArea.W
	origLen := len(ws.columns)
	col := &column{ws: ws, width: wsArea.W / uint32(origLen+1)}
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
		ws.columns = append([]*column{col}, ws.columns...)
	} else {
		ws.columns = append(ws.columns, col)
	}
	return col
}

func (ws *workspace) deleteColumn(col *column) {
	i := ws.findColumnIndex(func(c *column) bool { return c == col })
	if i < 0 {
		return
	}
	wsWidth := ws.output.workspaceArea().W
	// TODO: assign the widths proportional to the original width/totalWidth ratio
	// origLen = len(ws.columns)
	ws.columns = append(ws.columns[:i], ws.columns[i+1:]...)
	for _, c := range ws.columns {
		c.width = wsWidth / uint32(len(ws.columns))
	}
}

func (ws *workspace) findColumnIndex(predicate func(*column) bool) int {
	for i, col := range ws.columns {
		if predicate(col) {
			return i
		}
	}
	return -1
}

func (ws *workspace) updateTiling() {
	for _, col := range ws.columns {
		col.updateTiling()
	}
}

func (ws *workspace) fullArea() x11.Geom { return ws.output.workspaceArea() }

func (ws *workspace) area() x11.Geom {
	a := ws.fullArea()
	return x11.Geom{
		X: a.X + ws.config.gap,
		Y: a.Y + ws.config.gap,
		W: a.W - ws.config.gap*2,
		H: a.H - ws.config.gap*2,
	}
}

// singleFrame returns a single frame if there's only one in the workspace, nil otherwise
func (ws *workspace) singleFrame() *frame {
	if ws.countAllFrames() == 1 {
		return ws.columns[0].frames[0]
	}
	return nil
}

func (ws *workspace) countAllFrames() int {
	count := 0
	for _, col := range ws.columns {
		count += len(col.frames)
	}
	return count
}
