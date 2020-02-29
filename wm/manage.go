package wm

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/patrislav/marwind/client"
)

func (wm *WM) manageWindow(win xproto.Window) error {
	typ, err := wm.getWindowType(win)
	if err != nil {
		return fmt.Errorf("failed to get window type: %v", err)
	}
	mask := uint32(xproto.EventMaskStructureNotify | xproto.EventMaskEnterWindow | xproto.EventMaskPropertyChange)
	cookie := xproto.ChangeWindowAttributesChecked(wm.xc.X(), win, xproto.CwEventMask, []uint32{mask})
	if err := cookie.Check(); err != nil {
		return fmt.Errorf("failed to change window attributes: %v", err)
	}
	f, err := wm.createFrame(win, typ)
	if err != nil {
		return fmt.Errorf("failed to frame the window: %v", err)
	}
	switch f.cli.Type() {
	case client.TypeNormal:
		ws := wm.outputs[0].activeWs
		if err := ws.addFrame(f); err != nil {
			return fmt.Errorf("failed to add frame: %v", err)
		}
		if err := wm.renderWorkspace(ws); err != nil {
			return fmt.Errorf("failed to render workspace: %v", err)
		}
	case client.TypeDock:
		if err := wm.outputs[0].addDock(f); err != nil {
			return fmt.Errorf("failed to add dock: %v", err)
		}
		if err := wm.renderOutput(wm.outputs[0]); err != nil {
			return fmt.Errorf("failed to render output: %v", err)
		}
	}
	return nil
}

func (wm *WM) getWindowType(win xproto.Window) (client.Type, error) {
	typeAtom := wm.xc.Atom("_NET_WM_WINDOW_TYPE")
	dockTypeAtom := wm.xc.Atom("_NET_WM_WINDOW_TYPE_DOCK")
	normalTypeAtom := wm.xc.Atom("_NET_WM_WINDOW_TYPE_NORMAL")
	prop, err := xproto.GetProperty(wm.xc.X(), false, win, typeAtom, xproto.GetPropertyTypeAny, 0, 64).Reply()
	if err != nil {
		return client.TypeUnknown, err
	}
	if prop != nil {
		for v := prop.Value; len(v) >= 4; v = v[4:] {
			switch xproto.Atom(uint32(v[0]) | uint32(v[1])<<8 | uint32(v[2])<<16 | uint32(v[3])<<24) {
			case dockTypeAtom:
				return client.TypeDock, nil
			case normalTypeAtom:
				return client.TypeNormal, nil
			}
		}
	}
	return client.TypeNormal, nil
}
