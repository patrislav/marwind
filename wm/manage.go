package wm

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind/x11"
)

func (wm *WM) manageWindow(win xproto.Window) error {
	typ, err := getWindowType(win)
	if err != nil {
		return fmt.Errorf("failed to get window type: %v", err)
	}
	evtMask := []uint32{xproto.EventMaskStructureNotify | xproto.EventMaskEnterWindow}
	cookie := xproto.ChangeWindowAttributesChecked(x11.X, win, xproto.CwEventMask, evtMask)
	if err := cookie.Check(); err != nil {
		return fmt.Errorf("failed to change window attributes: %v", err)
	}
	f, err := wm.createFrame(win, typ)
	if err != nil {
		return fmt.Errorf("failed to frame the window: %v", err)
	}
	switch f.typ {
	case winTypeNormal:
		ws := wm.outputs[0].activeWs
		if err := ws.addFrame(f); err != nil {
			return fmt.Errorf("failed to add frame: %v", err)
		}
		if err := wm.renderWorkspace(ws); err != nil {
			return fmt.Errorf("failed to render workspace: %v", err)
		}
	case winTypeDock:
		if err := wm.outputs[0].addDock(f); err != nil {
			return fmt.Errorf("failed to add dock: %v", err)
		}
		if err := wm.renderOutput(wm.outputs[0]); err != nil {
			return fmt.Errorf("failed to render output: %v", err)
		}
	}
	return nil
}

func getWindowType(win xproto.Window) (winType, error) {
	typeAtom := x11.Atom("_NET_WM_WINDOW_TYPE")
	dockTypeAtom := x11.Atom("_NET_WM_WINDOW_TYPE_DOCK")
	normalTypeAtom := x11.Atom("_NET_WM_WINDOW_TYPE_NORMAL")
	prop, err := xproto.GetProperty(x11.X, false, win, typeAtom, xproto.GetPropertyTypeAny, 0, 64).Reply()
	if err != nil {
		return winTypeUnknown, err
	}
	if prop != nil {
		for v := prop.Value; len(v) >= 4; v = v[4:] {
			switch xproto.Atom(uint32(v[0]) | uint32(v[1])<<8 | uint32(v[2])<<16 | uint32(v[3])<<24) {
			case dockTypeAtom:
				return winTypeDock, nil
			case normalTypeAtom:
				return winTypeNormal, nil
			}
		}
	}
	return winTypeNormal, nil
}
