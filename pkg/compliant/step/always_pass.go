package step

import "bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"

type alwaysPass struct {
	order int
}

func NewAlwaysPass(order int) Step {
	return alwaysPass{
		order: order,
	}
}

func (s alwaysPass) Run(ctx context.Context) Result {
	return NewPassResult("always dumb pass step")
}

func (s alwaysPass) Order() int {
	return s.order
}
