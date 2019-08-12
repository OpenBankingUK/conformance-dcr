package step

import (
	"fmt"
	"net/http"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"
)

type getRequest struct {
	url        string
	contextVar string
	stepName   string
	httpClient *http.Client
}

func NewGetRequest(url, responseContextVar string, httpClient *http.Client) Step {
	return getRequest{
		url:        url,
		contextVar: responseContextVar,
		stepName:   fmt.Sprintf("GET request %s", url),
		httpClient: httpClient,
	}
}

func (s getRequest) Run(ctx context.Context) Result {
	r, err := s.httpClient.Get(s.url)
	if err != nil {
		return NewFailResult(s.stepName, err.Error())
	}

	ctx.SetResponse(s.contextVar, r)

	return NewPassResult(s.stepName)
}
