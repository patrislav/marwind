package wm

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind/x11"
)

func (wm *WM) setFocus(win xproto.Window, time xproto.Timestamp) error {
	frm := wm.findFrame(func(f *frame) bool { return f.client.window == win && f.typ == winTypeNormal })
	if frm == nil && win != x11.Screen.Root {
		return nil
	}
	wm.activeWin = win
	cookie := xproto.GetProperty(x11.X, false, win, x11.Atom("WM_PROTOCOLS"), xproto.GetPropertyTypeAny, 0, 64)
	prop, err := cookie.Reply()
	if err == nil && wm.takeFocusProp(prop, win, time) {
		return x11.SetActiveWindow(win)
	}
	err = xproto.SetInputFocusChecked(x11.X, xproto.InputFocusPointerRoot, win, time).Check()
	if err != nil {
		return err
	}
	return x11.SetActiveWindow(win)
}

func (wm *WM) removeFocus() error {
	return wm.setFocus(x11.Screen.Root, xproto.TimeCurrentTime)
}

func (wm *WM) takeFocusProp(prop *xproto.GetPropertyReply, win xproto.Window, time xproto.Timestamp) bool {
	for v := prop.Value; len(v) >= 4; v = v[4:] {
		switch xproto.Atom(uint32(v[0]) | uint32(v[1])<<8 | uint32(v[2])<<16 | uint32(v[3])<<24) {
		case x11.Atom("WM_TAKE_FOCUS"):
			_ = xproto.SendEventChecked(
				x11.X,
				false,
				win,
				xproto.EventMaskNoEvent,
				string(xproto.ClientMessageEvent{
					Format: 32,
					Window: win,
					Type:   x11.Atom("WM_PROTOCOLS"),
					Data: xproto.ClientMessageDataUnionData32New([]uint32{
						uint32(x11.Atom("WM_TAKE_FOCUS")),
						uint32(time),
						0,
						0,
						0,
					}),
				}.Bytes()),
			).Check()
			return true
		}
	}
	return false
}

func (wm *WM) warpPointerToFrame(f *frame) error {
	return x11.WarpPointer(f.geom.X+f.geom.W/2, f.geom.Y+f.geom.H/2)
}
