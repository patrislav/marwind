package manager

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind-wm/container"
	"github.com/patrislav/marwind-wm/keysym"
)

type Action struct {
	sym       xproto.Keysym
	modifiers int
	codes     []xproto.Keycode
	act       func() error
}

func initActions(m *Manager) []Action {
	mod1 := xproto.ModMask1
	shift := xproto.ModMaskShift
	// shift := 0
	actions := []Action{
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
	}
	for i, syms := range m.keymap {
		for _, sym := range syms {
			for c := range actions {
				if actions[c].sym == sym {
					fmt.Println(actions[c])
					actions[c].codes = append(actions[c].codes, xproto.Keycode(i))
				}
			}
		}
	}
	return actions
}

func handleRemoveWindow(m *Manager) error {
	if !m.ws.HasWindow(m.activeWin) {
		return nil
	}
	cookie := xproto.GetProperty(m.xc, false, m.activeWin, m.atoms.wmProtocols, xproto.GetPropertyTypeAny, 0, 64)
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
					m.xc,
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
		return xproto.DestroyWindowChecked(m.xc, m.activeWin).Check()
	}
	return nil
}

func handleMoveWindow(m *Manager, dir container.MoveDirection) error {
	err := m.ws.MoveWindow(m.activeWin, dir)
	if err != nil {
		return err
	}
	return m.renderWorkspace(m.ws)
}
