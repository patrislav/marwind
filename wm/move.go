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

func (wm *WM) switchWorkspace(id uint8) error {
	ws, err := wm.ensureWorkspace(id)
	if err != nil {
		return fmt.Errorf("failed to ensure workspace: %v", err)
	}
	if err := ws.output.switchWorkspace(ws); err != nil {
		return fmt.Errorf("output unable to switch workpace: %v", err)
	}
	if err := wm.renderWorkspace(ws); err != nil {
		return fmt.Errorf("wm.renderWorkspace: %w", err)
	}
	if err := wm.updateDesktopHints(); err != nil {
		return fmt.Errorf("failed to update desktop hints: %v", err)
	}
	if err := wm.removeFocus(); err != nil {
		return fmt.Errorf("failed to remove focus: %v", err)
	}

	// TODO: temporary solution! Focuses always the first window of the first column
	// Better approach: implement a window focus stack for each workspace, on switch focus the top-of-stack window
	if len(ws.columns) > 0 && len(ws.columns[0].frames) > 0 {
		win := ws.columns[0].frames[0].cli.Window()
		if err := wm.setFocus(win, xproto.TimeCurrentTime); err != nil {
			return fmt.Errorf("failed to set focus: %w", err)
		}
	}
	return nil
}

func (wm *WM) moveFrameToWorkspace(f *frame, wsID uint8) error {
	current := wm.outputs[0].activeWs
	next, err := wm.ensureWorkspace(wsID)
	if err != nil {
		return err
	}
	if next == current {
		return nil
	}
	if !current.deleteFrame(f) {
		return fmt.Errorf("frame not contained within workspace %d", wsID)
	}
	if err := next.addFrame(f); err != nil {
		return fmt.Errorf("failed to add the frame to the next workspace: %v", err)
	}
	if err := f.cli.Unmap(); err != nil {
		return fmt.Errorf("failed to unmap the frame: %v", err)
	}
	if err := wm.renderWorkspace(next); err != nil {
		return fmt.Errorf("failed to render next workspace: %v", err)
	}
	if err := wm.renderWorkspace(current); err != nil {
		return fmt.Errorf("failed to render previous workspace: %v", err)
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
