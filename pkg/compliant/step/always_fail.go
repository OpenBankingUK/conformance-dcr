package step

import "bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"

type alwaysFail struct{}

func NewAlwaysFail() Step {
	return alwaysFail{}
}

func (s alwaysFail) Run(ctx context.Context) Result {
	return NewFailResult("always dumb fail step", "failed because of reasons")
}
