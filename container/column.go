package container

type Column struct {
	ws     *Workspace
	width  uint32
	frames []*Frame
}

func (c *Column) AddFrame(frame *Frame, after *Frame) {
	frame.col = c
	wsHeight := c.ws.Rect().H
	if len(c.frames) > 0 {
		frame.height = wsHeight / uint32(len(c.frames)+1)
		remHeight := float32(wsHeight - frame.height)
		leftHeight := uint32(remHeight)
		for _, f := range c.frames {
			f.height = uint32((float32(f.height) / float32(wsHeight)) * remHeight)
			leftHeight -= f.height
		}
		if leftHeight != 0 {
			frame.height += leftHeight
		}
	} else {
		frame.height = wsHeight
	}
	c.frames = append(c.frames, frame)
}

func (c *Column) DeleteFrame(frame *Frame) {
	idx := c.findFrameIndex(func(f *Frame) bool { return f == frame })
	if idx < 0 {
		return
	}
	// origLen := len(c.frames)
	c.frames = append(c.frames[:idx], c.frames[idx+1:]...)
	c.UpdateTiling()
}

func (c *Column) UpdateTiling() {
	wsHeight := c.ws.Rect().H
	// TODO: assign the heights proportional to the original height/totalHeight ratio
	for _, f := range c.frames {
		f.height = wsHeight / uint32(len(c.frames))
	}
}

func (c *Column) Width() uint32    { return c.width }
func (c *Column) Frames() []*Frame { return c.frames }

func (c *Column) SetWidth(width uint32) {
	c.width = width
}

func (c *Column) findFrameIndex(predicate func(*Frame) bool) int {
	for i, f := range c.frames {
		if predicate(f) {
			return i
		}
	}
	return -1
}
