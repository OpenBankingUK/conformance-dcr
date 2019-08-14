package step

type alwaysPass struct{}

func NewAlwaysPass() Step {
	return alwaysPass{}
}

func (s alwaysPass) Run(ctx Context) Result {
	return NewPassResult("always dumb pass step")
}
