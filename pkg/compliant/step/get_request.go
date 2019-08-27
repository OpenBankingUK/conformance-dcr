package step

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type getRequest struct {
	url            string
	responseCtxKey string
	stepName       string
	httpClient     *http.Client
}

func NewGetRequest(url, responseContextVar string, httpClient *http.Client) Step {
	return getRequest{
		url:            url,
		responseCtxKey: responseContextVar,
		stepName:       fmt.Sprintf("GET request %s", url),
		httpClient:     httpClient,
	}
}

func (s getRequest) Run(ctx Context) Result {
	debug := NewDebug()

	debug.Logf("making get request to : %s", s.url)
	r, err := s.httpClient.Get(s.url)
	if err != nil {
		return NewFailResultWithDebug(s.stepName, err.Error(), debug)
	}
	debug.Logf("Response: %s", s.debugMessage(r))

	debug.Logf("setting response object in ctx var: %s", s.responseCtxKey)
	ctx.SetResponse(s.responseCtxKey, r)

	return NewPassResultWithDebug(s.stepName, debug)
}

func (s getRequest) debugMessage(r *http.Response) string {
	debug, err := httputil.DumpResponse(r, true)
	if err != nil {
		return fmt.Sprintf("could not dump response object: %s", err.Error())
	}
	return string(debug)
}
