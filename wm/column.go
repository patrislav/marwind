package wm

type column struct {
	ws     *workspace
	frames []*frame
	width  uint32
}

func (c *column) addFrame(frm *frame, after *frame) {
	frm.col = c
	wsHeight := c.ws.output.workspaceArea().H
	if len(c.frames) > 0 {
		frm.height = wsHeight / uint32(len(c.frames)+1)
		remHeight := float32(wsHeight - frm.height)
		leftHeight := uint32(remHeight)
		for _, f := range c.frames {
			f.height = uint32((float32(f.height) / float32(wsHeight)) * remHeight)
			leftHeight -= f.height
		}
		if leftHeight != 0 {
			frm.height += leftHeight
		}
	} else {
		frm.height = wsHeight
	}
	c.frames = append(c.frames, frm)
}

func (c *column) deleteFrame(frm *frame) {
	idx := c.findFrameIndex(func(f *frame) bool { return f == frm })
	if idx < 0 {
		return
	}
	c.frames = append(c.frames[:idx], c.frames[idx+1:]...)
	c.updateTiling()
}

func (c *column) updateTiling() {
	wsHeight := c.ws.output.workspaceArea().H
	// TODO: assign the heights proportional to the original height/totalHeight ratio
	for _, f := range c.frames {
		f.height = wsHeight / uint32(len(c.frames))
	}
}

func (c *column) findFrameIndex(predicate func(*frame) bool) int {
	for i, f := range c.frames {
		if predicate(f) {
			return i
		}
	}
	return -1
}
