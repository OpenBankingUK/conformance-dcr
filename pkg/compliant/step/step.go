package step

import "bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"

type Step interface {
	Run(ctx context.Context) Result
	Order() int
}

type Result struct {
	Name    string
	Pass    bool
	Message string
}

type Results []Result

func (r Results) Fail() bool {
	for _, result := range r {
		if !result.Pass {
			return true
		}
	}
	return false
}

func NewPassResult(name string) Result {
	return Result{Name: name, Pass: true}
}

func NewFailResult(name, msg string) Result {
	return Result{Name: name, Pass: false, Message: msg}
}
