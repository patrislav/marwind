package wm

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

func (wm *WM) switchWorkspace(id uint8) error {
	ws, err := wm.ensureWorkspace(id)
	if err != nil {
		return fmt.Errorf("failed to ensure workspace: %v", err)
	}
	if err := ws.output.switchWorkspace(ws); err != nil {
		return fmt.Errorf("output unable to switch workpace: %v", err)
	}
	if err := wm.removeFocus(); err != nil {
		return fmt.Errorf("failed to remove focus: %v", err)
	}
	if err := wm.updateDesktopHints(); err != nil {
		return fmt.Errorf("failed to update desktop hints: %v", err)
	}
	return nil
}

// ensureWorkspace looks up a workspace by ID, adding it to the current output if needed
func (wm *WM) ensureWorkspace(id uint8) (*workspace, error) {
	var nextWs *workspace
	for _, ws := range wm.workspaces {
		if ws.id == id {
			nextWs = ws
			break
		}
	}
	if nextWs == nil {
		return nil, fmt.Errorf("no workspace with ID %d", id)
	}
	switch {
	case nextWs.output == nil:
		if err := wm.outputs[0].addWorkspace(nextWs); err != nil {
			return nil, err
		}
	case nextWs.output != wm.outputs[0]:
		return nil, fmt.Errorf("multiple outputs not supported yet")
	}
	return nextWs, nil
}
