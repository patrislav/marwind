package wm

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind/client"
)

func (wm *WM) setFocus(win xproto.Window, time xproto.Timestamp) error {
	frm := wm.findFrame(func(f *frame) bool { return f.cli.Window() == win && f.cli.Type() == client.TypeNormal })
	if frm == nil && win != wm.xc.GetRootWindow() {
		return nil
	}
	wm.activeWin = win
	cookie := xproto.GetProperty(wm.xc.X(), false, win, wm.xc.Atom("WM_PROTOCOLS"), xproto.GetPropertyTypeAny, 0, 64)
	prop, err := cookie.Reply()
	if err == nil && wm.takeFocusProp(prop, win, time) {
		return wm.xc.SetActiveWindow(win)
	}
	err = xproto.SetInputFocusChecked(wm.xc.X(), xproto.InputFocusPointerRoot, win, time).Check()
	if err != nil {
		return err
	}
	return wm.xc.SetActiveWindow(win)
}

func (wm *WM) removeFocus() error {
	return wm.setFocus(wm.xc.GetRootWindow(), xproto.TimeCurrentTime)
}

func (wm *WM) takeFocusProp(prop *xproto.GetPropertyReply, win xproto.Window, time xproto.Timestamp) bool {
	for v := prop.Value; len(v) >= 4; v = v[4:] {
		switch xproto.Atom(uint32(v[0]) | uint32(v[1])<<8 | uint32(v[2])<<16 | uint32(v[3])<<24) {
		case wm.xc.Atom("WM_TAKE_FOCUS"):
			_ = xproto.SendEventChecked(
				wm.xc.X(),
				false,
				win,
				xproto.EventMaskNoEvent,
				string(xproto.ClientMessageEvent{
					Format: 32,
					Window: win,
					Type:   wm.xc.Atom("WM_PROTOCOLS"),
					Data: xproto.ClientMessageDataUnionData32New([]uint32{
						uint32(wm.xc.Atom("WM_TAKE_FOCUS")),
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
	geom := f.cli.Geom()
	return wm.xc.WarpPointer(geom.X+int16(geom.W/2), geom.Y+int16(geom.H/2))
}
