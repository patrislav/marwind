package container

type Output struct {
	rect       Rect
	workspaces []*Workspace
	currentWs  *Workspace
}

func NewOutput(rect Rect) *Output {
	return &Output{rect: rect}
}

func (o *Output) AddWorkspace(ws *Workspace) {
	ws.setOutput(o)
	o.workspaces = append(o.workspaces, ws)
}
