package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/http"
	"fmt"
)

type assertStatusCode struct {
	code               int
	responseContextVar string
	stepName           string
}

func NewAssertStatus(code int, responseContextVar string) Step {
	return assertStatusCode{
		code:               code,
		responseContextVar: responseContextVar,
		stepName:           fmt.Sprintf("Assert status code %d", code),
	}
}

func (a assertStatusCode) Run(ctx Context) Result {
	debug := NewDebug()

	debug.Logf("get response object from ctx var: %s", a.responseContextVar)
	r, err := ctx.GetResponse(a.responseContextVar)
	if err != nil {
		return NewFailResult(a.stepName, fmt.Sprintf("getting response object from context: %s", err.Error()))
	}

	if r.StatusCode != a.code {
		debug.Log(http.DebugResponse(r))
		return NewFailResultWithDebug(
			a.stepName,
			fmt.Sprintf("Expecting status code %d but got %d", a.code, r.StatusCode),
			debug,
		)
	}

	return NewPassResultWithDebug(a.stepName, debug)
}
