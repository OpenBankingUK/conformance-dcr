package step

type alwaysFail struct{}

func NewAlwaysFail() Step {
	return alwaysFail{}
}

func (s alwaysFail) Run(ctx Context) Result {
	debug := NewDebug()
	debug.Log("always fail step")
	return NewFailResultWithDebug("always dumb fail step", "failed because of reasons", debug)
}
