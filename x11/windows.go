package x11

import (
	"fmt"
	"time"

	"github.com/BurntSushi/xgb/xproto"
)

func GracefullyDestroyWindow(win xproto.Window) error {
	protos, err := getProps32(win, "WM_PROTOCOLS")
	if err != nil {
		return fmt.Errorf("could not close window: %v", err)
	}
	for _, p := range protos {
		if xproto.Atom(p) == Atom("WM_DELETE_WINDOW") {
			t := time.Now().Unix()
			return xproto.SendEventChecked(
				X,
				false,
				win,
				xproto.EventMaskNoEvent,
				string(xproto.ClientMessageEvent{
					Format: 32,
					Window: win,
					Type:   Atom("WM_PROTOCOLS"),
					Data: xproto.ClientMessageDataUnionData32New([]uint32{
						uint32(Atom("WM_DELETE_WINDOW")),
						uint32(t),
						0,
						0,
						0,
					}),
				}.Bytes()),
			).Check()
		}
	}
	// The window does not follow ICCCM - just destroy it
	return xproto.DestroyWindowChecked(X, win).Check()
}
