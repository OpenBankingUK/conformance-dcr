package step

type alwaysPass struct{}

func NewAlwaysPass() Step {
	return alwaysPass{}
}

func (s alwaysPass) Run(ctx Context) Result {
	debug := NewDebug()
	debug.Log("always fail step")
	return NewPassResultWithDebug("always dumb pass step", debug)
}
