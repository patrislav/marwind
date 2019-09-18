package manager

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind/container"
	"github.com/patrislav/marwind/keysym"
	"github.com/patrislav/marwind/x11"
)

type Action struct {
	sym       xproto.Keysym
	modifiers int
	codes     []xproto.Keycode
	act       func() error
}

func initActions(m *Manager) []*Action {
	mod1 := xproto.ModMask1
	shift := xproto.ModMaskShift
	actions := []*Action{
		{
			sym:       keysym.XK_q,
			modifiers: mod1,
			act: func() error {
				return handleRemoveWindow(m)
			},
		},
		{
			sym:       keysym.XK_d,
			modifiers: mod1,
			act: func() error {
				cmd := exec.Command("rofi", "-show", "drun")
				err := cmd.Run()
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			sym:       keysym.XK_h,
			modifiers: mod1 | shift,
			act:       func() error { return handleMoveWindow(m, container.MoveLeft) },
		},
		{
			sym:       keysym.XK_j,
			modifiers: mod1 | shift,
			act:       func() error { return handleMoveWindow(m, container.MoveDown) },
		},
		{
			sym:       keysym.XK_k,
			modifiers: mod1 | shift,
			act:       func() error { return handleMoveWindow(m, container.MoveUp) },
		},
		{
			sym:       keysym.XK_l,
			modifiers: mod1 | shift,
			act:       func() error { return handleMoveWindow(m, container.MoveRight) },
		},
		{
			sym:       keysym.XK_y,
			modifiers: mod1 | shift,
			act:       func() error { return handleResizeWindow(m, container.ResizeHoriz, -5) },
		},
		{
			sym:       keysym.XK_u,
			modifiers: mod1 | shift,
			act:       func() error { return handleResizeWindow(m, container.ResizeVert, 5) },
		},
		{
			sym:       keysym.XK_i,
			modifiers: mod1 | shift,
			act:       func() error { return handleResizeWindow(m, container.ResizeVert, -5) },
		},
		{
			sym:       keysym.XK_o,
			modifiers: mod1 | shift,
			act:       func() error { return handleResizeWindow(m, container.ResizeHoriz, 5) },
		},
	}
	actions = appendWorkspaceActions(m, actions, mod1)
	for i, syms := range m.keymap {
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

func appendWorkspaceActions(m *Manager, actions []*Action, switchMod int) []*Action {
	for i := 0; i < maxWorkspaces; i++ {
		var sym xproto.Keysym
		if i == 9 {
			sym = keysym.XK_0
		} else {
			sym = xproto.Keysym(keysym.XK_1 + i)
		}
		wsID := i
		actions = append(actions, &Action{
			sym:       sym,
			modifiers: switchMod,
			act: func() error {
				fmt.Println("switching", wsID)
				return m.switchWorkspace(m.outputs[0], uint8(wsID))
			},
		})
	}
	return actions
}

func handleRemoveWindow(m *Manager) error {
	if !m.outputs[0].CurrentWorkspace().HasWindow(m.activeWin) {
		return nil
	}
	cookie := xproto.GetProperty(x11.X, false, m.activeWin, m.atoms.wmProtocols, xproto.GetPropertyTypeAny, 0, 64)
	prop, err := cookie.Reply()
	if err != nil {
		log.Println("error when getting property", err)
		return err
	}
	if prop != nil {
		for v := prop.Value; len(v) >= 4; v = v[4:] {
			switch xproto.Atom(uint32(v[0]) | uint32(v[1])<<8 | uint32(v[2])<<16 | uint32(v[3])<<24) {
			case m.atoms.wmDeleteWindow:
				t := time.Now().Unix()
				return xproto.SendEventChecked(
					x11.X,
					false,
					m.activeWin,
					xproto.EventMaskNoEvent,
					string(xproto.ClientMessageEvent{
						Format: 32,
						Window: m.activeWin,
						Type:   m.atoms.wmProtocols,
						Data: xproto.ClientMessageDataUnionData32New([]uint32{
							uint32(m.atoms.wmDeleteWindow),
							uint32(t),
							0,
							0,
							0,
						}),
					}.Bytes())).Check()
			}
		}
	}
	// There were no properties which means window doesn't follow ICCCM. Just destroy it
	if m.activeWin != 0 {
		return xproto.DestroyWindowChecked(x11.X, m.activeWin).Check()
	}
	return nil
}

func handleMoveWindow(m *Manager, dir container.MoveDirection) error {
	ws := m.outputs[0].CurrentWorkspace()
	err := ws.MoveWindow(m.activeWin, dir)
	if err != nil {
		return err
	}
	err = m.renderWorkspace(ws)
	if err != nil {
		return err
	}
	frame := m.findFrame(func(f *container.Frame) bool { return f.Window() == m.activeWin })
	if frame == nil {
		return fmt.Errorf("could not find frame with window %v", m.activeWin)
	}
	return m.warpPointerToFrame(frame)
}

func handleResizeWindow(m *Manager, dir container.ResizeDirection, pct int) error {
	ws := m.outputs[0].CurrentWorkspace()
	err := ws.ResizeWindow(m.activeWin, dir, pct)
	if err != nil {
		return err
	}
	err = m.renderWorkspace(ws)
	if err != nil {
		return err
	}
	return nil
}
