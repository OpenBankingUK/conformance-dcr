package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"
	"fmt"
	"net/http"
)

type getRequest struct {
	url        string
	contextVar string
	order      int
	stepName   string
}

func NewGetRequest(order int, url, responseContextVar string) Step {
	return getRequest{
		url:        url,
		contextVar: responseContextVar,
		order:      order,
		stepName:   fmt.Sprintf("GET request %s", url),
	}
}

func (s getRequest) Run(ctx context.Context) Result {
	c := http.Client{}

	r, err := c.Get(s.url)
	if err != nil {
		return NewFailResult(s.stepName, err.Error())
	}

	ctx.SetResponse(s.contextVar, r)

	return NewPassResult(s.stepName)
}

func (s getRequest) Order() int {
	return s.order
}
