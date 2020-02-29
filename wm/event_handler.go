package wm

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"
)

type eventHandler struct {
	wm *WM
}

func (h eventHandler) eventLoop() {
	for {
		xev, err := h.wm.xc.X().WaitForEvent()
		if err != nil {
			log.Println(err)
			continue
		}
		switch e := xev.(type) {
		case xproto.KeyPressEvent:
			h.keyPress(e)
		case xproto.EnterNotifyEvent:
			h.enterNotify(e)
		case xproto.ConfigureRequestEvent:
			h.configureRequest(e)
		case xproto.MapNotifyEvent:
			h.mapNotify(e)
		case xproto.MapRequestEvent:
			h.mapRequest(e)
		case xproto.UnmapNotifyEvent:
			h.unmapNotify(e)
		case xproto.DestroyNotifyEvent:
			h.destroyNotify(e)
		case xproto.PropertyNotifyEvent:
			h.propertyNotify(e)
		case xproto.ClientMessageEvent:
			h.clientMessage(e)
		case xproto.ExposeEvent:
			h.expose(e)
		}
	}
}

func (h eventHandler) keyPress(e xproto.KeyPressEvent) {
	if err := h.wm.handleKeyPressEvent(e); err != nil {
		log.Println(err)
	}
}

func (h eventHandler) enterNotify(e xproto.EnterNotifyEvent) {
	f := h.wm.findFrame(func(frm *frame) bool { return frm.cli.Window() == e.Event })
	if f != nil {
		if err := h.wm.setFocus(e.Event, e.Time); err != nil {
			log.Println("Failed to set focus:", err)
		}
	}
}

func (h eventHandler) configureRequest(e xproto.ConfigureRequestEvent) {
	if err := h.wm.handleConfigureRequest(e); err != nil {
		log.Println("Failed to configure window:", err)
	}
}

func (h eventHandler) mapNotify(e xproto.MapNotifyEvent) {
	f := h.wm.findFrame(func(frm *frame) bool { return frm.cli.Window() == e.Window })
	if f != nil {
		if err := h.wm.configureNotify(f); err != nil {
			log.Printf("Failed to send ConfigureNotify event to %d: %v\n", e.Window, err)
		}
	}
}

func (h eventHandler) mapRequest(e xproto.MapRequestEvent) {
	f := h.wm.findFrame(func(frm *frame) bool { return frm.cli.Window() == e.Window })
	if f != nil {
		log.Printf("Skipping MapRequest of an already mapped window %d\n", e.Window)
		return
	}
	if attr, err := xproto.GetWindowAttributes(h.wm.xc.X(), e.Window).Reply(); err != nil || !attr.OverrideRedirect {
		if err := h.wm.manageWindow(e.Window); err != nil {
			log.Println("Failed to manage a window:", err)
		}
	}
	if err := h.wm.updateDesktopHints(); err != nil {
		log.Printf("Failed to update desktop hints: %v", err)
	}
}

func (h eventHandler) unmapNotify(e xproto.UnmapNotifyEvent) {
	f := h.wm.findFrame(func(frm *frame) bool { return frm.cli.Window() == e.Window })
	if f != nil {
		if err := f.cli.OnUnmap(); err != nil {
			log.Println("Failed to unmap frame's parent:", err)
			return
		}
	}
}

func (h eventHandler) destroyNotify(e xproto.DestroyNotifyEvent) {
	f := h.wm.findFrame(func(frm *frame) bool { return frm.cli.Window() == e.Window })
	if f != nil {
		if err := f.cli.OnDestroy(); err != nil {
			log.Println("Failed to destroy frame's parent:", err)
			return
		}
		if err := h.wm.deleteFrame(f); err != nil {
			log.Println("Failed to delete the frame:", err)
		}
		if err := h.wm.updateDesktopHints(); err != nil {
			log.Printf("Failed to update desktop hints: %v", err)
		}
	}
}

func (h eventHandler) propertyNotify(e xproto.PropertyNotifyEvent) {
	f := h.wm.findFrame(func(frm *frame) bool { return frm.cli.Window() == e.Window })
	if f != nil {
		f.cli.OnProperty(e.Atom)
	}
}

func (h eventHandler) clientMessage(e xproto.ClientMessageEvent) {
	switch e.Type {
	case h.wm.xc.Atom("_NET_CURRENT_DESKTOP"):
		out := h.wm.outputs[0]
		index := int(e.Data.Data32[0])
		if index < len(out.workspaces) {
			ws := out.workspaces[index]
			if err := h.wm.switchWorkspace(ws.id); err != nil {
				log.Printf("Failed to switch workspace: %v", err)
			}
		}
	}
}

func (h eventHandler) expose(e xproto.ExposeEvent) {
	f := h.wm.findFrame(func(frm *frame) bool {
		return frm.cli.Parent() == e.Window || frm.cli.Window() == e.Window
	})
	if f != nil {
		if err := f.cli.Draw(); err != nil {
			log.Println("Failed to draw client:", err)
		}
	}
}
