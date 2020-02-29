package client

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"
)

type Geom struct {
	X, Y int16
	W, H uint16
}

type Type uint8

const (
	TypeUnknown Type = iota
	TypeNormal
	TypeDock
)

type Client struct {
	x11    x11
	window xproto.Window
	parent xproto.Window
	mapped bool

	geom Geom
	cfg  *Config
	typ  Type

	title string
}

func New(x11 x11, cfg *Config, window xproto.Window, typ Type) (*Client, error) {
	c := &Client{x11: x11, cfg: cfg, window: window, typ: typ}

	if typ == TypeNormal {
		parent, err := c.createParent()
		if err != nil {
			return nil, fmt.Errorf("failed to create parent: %w", err)
		}
		if err := c.reparent(parent); err != nil {
			return nil, err
		}
		c.updateTitleProperty()
	}

	return c, nil
}

func (c *Client) Type() Type            { return c.typ }
func (c *Client) Window() xproto.Window { return c.window }
func (c *Client) Parent() xproto.Window { return c.parent }
func (c *Client) Geom() Geom            { return c.geom }
func (c *Client) Mapped() bool          { return c.mapped }
func (c *Client) SetGeom(geom Geom)     { c.geom = geom }

func (c *Client) Draw() error {
	return c.drawTitlebar()
}

// Update compares the desired state of the client against the actual state and executes updates
// aimed at reaching the desired state
func (c *Client) Update() error {
	return nil
}

// Map causes both the client window and the frame (parent) to be mapped
func (c *Client) Map() error {
	if c.parent != 0 {
		if err := c.x11.MapWindow(c.parent); err != nil {
			return fmt.Errorf("could not map parent: %w", err)
		}
	}
	if err := c.x11.MapWindow(c.window); err != nil {
		return fmt.Errorf("could not map window: %w", err)
	}
	c.mapped = true
	return nil
}

// Unmap causes the client window to be unmapped. This in turn sends the UnmapNotify event
// that is then handled by (*Client).OnUnmap
func (c *Client) Unmap() error {
	if err := c.x11.UnmapWindow(c.window); err != nil {
		return fmt.Errorf("could not unmap window: %w", err)
	}
	return nil
}

// OnDestroy is called when the WM receives the DestroyNotify event
func (c *Client) OnDestroy() error {
	if c.parent != 0 {
		if err := c.x11.DestroyWindow(c.parent); err != nil {
			return fmt.Errorf("could not destroy parent: %w", err)
		}
	}
	return nil
}

// OnUnmap is called when the WM receives the UnmapNotify event (e.g. when the client window
// is closed by user action or when requested by the program itself)
func (c *Client) OnUnmap() error {
	if !c.mapped {
		return nil
	}
	if c.parent != 0 {
		if err := c.x11.UnmapWindow(c.parent); err != nil {
			return fmt.Errorf("could not unmap parent: %w", err)
		}
	}
	c.mapped = false
	return nil
}

func (c *Client) OnProperty(atom xproto.Atom) {
	switch atom {
	case c.x11.Atom("_NET_WM_NAME"):
		c.updateTitleProperty()
	}
}

// createParent generates an X window and sets it up so that it can be used for reparenting
func (c *Client) createParent() (xproto.Window, error) {
	return c.x11.CreateWindow(c.x11.GetRootWindow(),
		0, 0, 1, 1, 0, xproto.WindowClassInputOutput,
		xproto.CwBackPixel|xproto.CwOverrideRedirect|xproto.CwEventMask,
		[]uint32{
			0xffa1d1cf,
			1,
			xproto.EventMaskSubstructureRedirect |
				xproto.EventMaskExposure |
				xproto.EventMaskButtonPress |
				xproto.EventMaskButtonRelease |
				xproto.EventMaskFocusChange,
		},
	)
}

func (c *Client) reparent(parent xproto.Window) error {
	if err := c.x11.ReparentWindow(c.window, parent, 0, 0); err != nil {
		return fmt.Errorf("could not reparent window: %w", err)
	}
	c.parent = parent
	return nil
}

func (c *Client) updateTitleProperty() {
	if v, err := c.x11.GetWindowTitle(c.window); err == nil {
		c.title = v
		c.drawTitlebar()
	}
}
