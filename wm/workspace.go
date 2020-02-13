package wm

import (
	"github.com/patrislav/marwind/client"
)

type workspaceConfig struct {
	gap uint16
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
		return f.cli.Map()
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

// resizeFrame changes the size of the frame by the given percent
func (ws *workspace) resizeFrame(f *frame, dir ResizeDirection, pct int) error {
	switch dir {
	case ResizeHoriz:
		if len(ws.columns) < 2 {
			return nil
		}
		min := uint16(float32(ws.area().W) * 0.1)
		dwFull := int(float32(ws.area().W) * (float32(pct) / 100))
		if uint16(int(f.col.width)+dwFull) < min {
			return nil
		}
		dwPart := dwFull/len(ws.columns) - 1
		dwFinal := 0
		for _, col := range ws.columns {
			if col != f.col {
				next := uint16(int(col.width) - dwPart)
				if next >= min {
					col.width = next
					dwFinal += dwPart
				}
			}
		}
		f.col.width = uint16(int(f.col.width) + dwFinal)
	case ResizeVert:
		col := f.col
		if len(col.frames) < 2 {
			return nil
		}
		min := uint16(float32(ws.area().H) * 0.1)
		dhFull := int(float32(ws.area().H) * (float32(pct) / 100))
		if uint16(int(f.height)+dhFull) < min {
			return nil
		}
		dhPart := dhFull/len(col.frames) - 1
		dhFinal := 0
		for _, other := range col.frames {
			if f != other {
				next := uint16(int(other.height) - dhPart)
				if next >= min {
					other.height = next
					dhFinal += dhPart
				}
			}
		}
		f.height = uint16(int(f.height) + dhFinal)
	}
	return nil
}

// show maps all the frames of the workspace
func (ws *workspace) show() error {
	var err error
	for _, col := range ws.columns {
		for _, f := range col.frames {
			if e := f.cli.Map(); e != nil {
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
			if e := f.cli.Unmap(); e != nil {
				err = e
			}
		}
	}
	return err
}

// createColumn creates a new empty column either at the start (if the start argument is true)
// or the end of the workspace area.
func (ws *workspace) createColumn(start bool) *column {
	wsWidth := ws.area().W
	origLen := len(ws.columns)
	col := &column{ws: ws, width: ws.area().W / uint16(origLen+1)}
	if origLen > 0 {
		col.width = wsWidth / uint16(origLen+1)
		remWidth := float32(wsWidth - col.width)
		leftWidth := uint16(remWidth)
		for _, c := range ws.columns {
			c.width = uint16((float32(c.width) / float32(wsWidth)) * remWidth)
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
		c.width = wsWidth / uint16(len(ws.columns))
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

func (ws *workspace) fullArea() client.Geom { return ws.output.workspaceArea() }

func (ws *workspace) area() client.Geom {
	a := ws.fullArea()
	return client.Geom{
		X: a.X + int16(ws.config.gap),
		Y: a.Y + int16(ws.config.gap),
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
