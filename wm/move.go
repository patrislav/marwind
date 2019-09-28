package wm

import (
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

func (wm *WM) moveWindow(win xproto.Window, dir MoveDirection) error {
	f := wm.findFrame(func(f *frame) bool { return f.client.window == win })
	if f == nil {
		return nil
	}
	ws := f.col.ws
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
