package step

type alwaysFail struct{}

func NewAlwaysFail() Step {
	return alwaysFail{}
}

func (s alwaysFail) Run(ctx Context) Result {
	return NewFailResult("always dumb fail step", "failed because of reasons")
}
