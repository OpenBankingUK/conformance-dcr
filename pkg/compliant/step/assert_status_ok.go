package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"
	"fmt"
	"net/http"
)

type assertStatusOk struct {
	responseContextVar string
	order              int
	stepName           string
}

func NewAssertStatusOk(order int, responseContextVar string) Step {
	return assertStatusOk{
		responseContextVar: responseContextVar,
		order:              order,
		stepName:           "Status Code 200",
	}
}

func (a assertStatusOk) Run(ctx context.Context) Result {
	response, err := ctx.GetResponse(a.responseContextVar)
	if err != nil {
		return NewFailResult(a.stepName, fmt.Sprintf("getting response object from context: %s", err.Error()))
	}

	if response.StatusCode != http.StatusOK {
		return NewFailResult(a.stepName, fmt.Sprintf("status received %d", response.StatusCode))
	}

	return NewPassResult(a.stepName)
}

func (a assertStatusOk) Order() int {
	return a.order
}
