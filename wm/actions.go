package wm

import (
	"log"
	"os"
	"os/exec"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind/keysym"
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
			sym:       keysym.XKq,
			modifiers: mod | shift,
			act: func() error {
				return handleRemoveWindow(wm)
			},
		},
		{
			sym:       keysym.XKt,
			modifiers: mod | shift | xproto.ModMask1,
			act: func() error {
				os.Exit(1)
				return nil
			},
		},
		{
			sym:       keysym.XKd,
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
			sym:       keysym.XKReturn,
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
			sym:       keysym.XKh,
			modifiers: mod | shift,
			act:       func() error { return handleMoveWindow(wm, MoveLeft) },
		},
		{
			sym:       keysym.XKj,
			modifiers: mod | shift,
			act:       func() error { return handleMoveWindow(wm, MoveDown) },
		},
		{
			sym:       keysym.XKk,
			modifiers: mod | shift,
			act:       func() error { return handleMoveWindow(wm, MoveUp) },
		},
		{
			sym:       keysym.XKl,
			modifiers: mod | shift,
			act:       func() error { return handleMoveWindow(wm, MoveRight) },
		},
		{
			sym:       keysym.XKy,
			modifiers: mod | shift,
			act:       func() error { return handleResizeWindow(wm, ResizeHoriz, -5) },
		},
		{
			sym:       keysym.XKu,
			modifiers: mod | shift,
			act:       func() error { return handleResizeWindow(wm, ResizeVert, 5) },
		},
		{
			sym:       keysym.XKi,
			modifiers: mod | shift,
			act:       func() error { return handleResizeWindow(wm, ResizeVert, -5) },
		},
		{
			sym:       keysym.XKo,
			modifiers: mod | shift,
			act:       func() error { return handleResizeWindow(wm, ResizeHoriz, 5) },
		},
	}
	actions = appendWorkspaceActions(wm, actions, mod, mod|shift)

	for sym, command := range wm.config.Keybindings {
		cmd := command
		actions = append(actions, &action{
			sym: sym,
			act: func() error {
				cmd := exec.Command(wm.config.Shell, "-c", cmd)
				go func() {
					if err := cmd.Run(); err != nil {
						log.Printf("Failed to run command (%s): %v\n", cmd, err)
					}
				}()
				return nil
			},
		})
	}

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
			sym = keysym.XK0
		} else {
			sym = xproto.Keysym(keysym.XK1 + i)
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
	frm := wm.findFrame(func(f *frame) bool { return f.cli.Window() == wm.activeWin })
	if frm == nil {
		log.Printf("WARNING: handleRemoveWindow: could not find frame with window %d\n", wm.activeWin)
		return nil
	}
	return wm.xc.GracefullyDestroyWindow(frm.cli.Window())
}

func handleMoveWindow(wm *WM, dir MoveDirection) error {
	frm := wm.findFrame(func(f *frame) bool { return f.cli.Window() == wm.activeWin })
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
	frm := wm.findFrame(func(f *frame) bool { return f.cli.Window() == wm.activeWin })
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
	frm := wm.findFrame(func(f *frame) bool { return f.cli.Window() == wm.activeWin })
	if frm == nil {
		log.Printf("WARNING: handleMoveWindowToWorkspace: could not find frame with window %d\n", wm.activeWin)
		return nil
	}
	if err := wm.moveFrameToWorkspace(frm, wsID); err != nil {
		return err
	}
	return nil
}
