package wm

import (
	"log"
	"os/exec"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind/keysym"
	"github.com/patrislav/marwind/x11"
)

type action struct {
	sym       xproto.Keysym
	modifiers int
	codes     []xproto.Keycode
	act       func() error
}

func initActions(wm *WM) []*action {
	mod := xproto.ModMask4
	shift := xproto.ModMaskShift
	actions := []*action{
		{
			sym:       keysym.XK_q,
			modifiers: mod | shift,
			act: func() error {
				return handleRemoveWindow(wm)
			},
		},
		{
			sym:       keysym.XK_d,
			modifiers: mod,
			act: func() error {
				cmd := exec.Command(wm.config.Shell, "-c", wm.config.LauncherCommand)
				go func() {
					if err := cmd.Run(); err != nil {
						log.Println("Failed to open launcher:", err)
					}
				}()
				return nil
			},
		},
		{
			sym:       keysym.XK_Return,
			modifiers: mod | shift,
			act: func() error {
				cmd := exec.Command(wm.config.Shell, "-c", wm.config.TerminalCommand)
				go func() {
					if err := cmd.Run(); err != nil {
						log.Println("Failed to open terminal:", err)
					}
				}()
				return nil
			},
		},
		{
			sym:       keysym.XK_h,
			modifiers: mod | shift,
			act:       func() error { return handleMoveWindow(wm, MoveLeft) },
		},
		{
			sym:       keysym.XK_j,
			modifiers: mod | shift,
			act:       func() error { return handleMoveWindow(wm, MoveDown) },
		},
		{
			sym:       keysym.XK_k,
			modifiers: mod | shift,
			act:       func() error { return handleMoveWindow(wm, MoveUp) },
		},
		{
			sym:       keysym.XK_l,
			modifiers: mod | shift,
			act:       func() error { return handleMoveWindow(wm, MoveRight) },
		},
		{
			sym:       keysym.XK_y,
			modifiers: mod | shift,
			act:       func() error { return handleResizeWindow(wm, ResizeHoriz, -5) },
		},
		{
			sym:       keysym.XK_u,
			modifiers: mod | shift,
			act:       func() error { return handleResizeWindow(wm, ResizeVert, 5) },
		},
		{
			sym:       keysym.XK_i,
			modifiers: mod | shift,
			act:       func() error { return handleResizeWindow(wm, ResizeVert, -5) },
		},
		{
			sym:       keysym.XK_o,
			modifiers: mod | shift,
			act:       func() error { return handleResizeWindow(wm, ResizeHoriz, 5) },
		},
	}
	actions = appendWorkspaceActions(wm, actions, mod, mod|shift)
	for i, syms := range wm.keymap {
		for _, sym := range syms {
			for c := range actions {
				if actions[c].sym == sym {
					actions[c].codes = append(actions[c].codes, xproto.Keycode(i))
				}
			}
		}
	}
	return actions
}

func appendWorkspaceActions(wm *WM, actions []*action, switchMod int, moveMod int) []*action {
	for i := 0; i < maxWorkspaces; i++ {
		var sym xproto.Keysym
		if i == 9 {
			sym = keysym.XK_0
		} else {
			sym = xproto.Keysym(keysym.XK_1 + i)
		}
		wsID := i
		actions = append(actions, &action{
			sym:       sym,
			modifiers: switchMod,
			act: func() error {
				return handleSwitchWorkspace(wm, uint8(wsID))
			},
		}, &action{
			sym:       sym,
			modifiers: moveMod,
			act: func() error {
				return handleMoveWindowToWorkspace(wm, uint8(wsID))
			},
		})
	}
	return actions
}

func handleRemoveWindow(wm *WM) error {
	frm := wm.findFrame(func(f *frame) bool { return f.client.window == wm.activeWin })
	if frm == nil {
		log.Printf("WARNING: handleRemoveWindow: could not find frame with window %d\n", wm.activeWin)
		return nil
	}
	return x11.GracefullyDestroyWindow(frm.client.window)
}

func handleMoveWindow(wm *WM, dir MoveDirection) error {
	frm := wm.findFrame(func(f *frame) bool { return f.client.window == wm.activeWin })
	if frm == nil {
		log.Printf("WARNING: handleMoveWindow: could not find frame with window %d\n", wm.activeWin)
		return nil
	}
	if err := frm.workspace().moveFrame(frm, dir); err != nil {
		return err
	}
	if err := wm.renderWorkspace(frm.workspace()); err != nil {
		return err
	}
	return wm.warpPointerToFrame(frm)
}

func handleResizeWindow(wm *WM, dir ResizeDirection, pct int) error {
	frm := wm.findFrame(func(f *frame) bool { return f.client.window == wm.activeWin })
	if frm == nil {
		log.Printf("WARNING: handleResizeWindow: could not find frame with window %d\n", wm.activeWin)
		return nil
	}
	if err := frm.workspace().resizeFrame(frm, dir, pct); err != nil {
		return err
	}
	if err := wm.renderWorkspace(frm.workspace()); err != nil {
		return err
	}
	return wm.warpPointerToFrame(frm)
}

func handleSwitchWorkspace(wm *WM, wsID uint8) error {
	return wm.switchWorkspace(wsID)
}

func handleMoveWindowToWorkspace(wm *WM, wsID uint8) error {
	frm := wm.findFrame(func(f *frame) bool { return f.client.window == wm.activeWin })
	if frm == nil {
		log.Printf("WARNING: handleMoveWindowToWorkspace: could not find frame with window %d\n", wm.activeWin)
		return nil
	}
	if err := wm.moveFrameToWorkspace(frm, wsID); err != nil {
		return err
	}
	return nil
}
