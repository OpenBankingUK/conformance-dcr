package step

import "bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"

type alwaysPass struct{}

func NewAlwaysPass() Step {
	return alwaysPass{}
}

func (s alwaysPass) Run(ctx context.Context) Result {
	return NewPassResult("always dumb pass step")
}
