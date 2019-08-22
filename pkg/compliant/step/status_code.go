package step

import (
	"fmt"
	"net/http"
	"net/http/httputil"
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
	response, err := ctx.GetResponse(a.responseContextVar)
	if err != nil {
		return NewFailResult(a.stepName, fmt.Sprintf("getting response object from context: %s", err.Error()))
	}

	if response.StatusCode != a.code {
		return NewFailResult(a.stepName, a.debugMessage(response))
	}

	return NewPassResult(a.stepName)
}

func (a assertStatusCode) debugMessage(r *http.Response) string {
	debug, err := httputil.DumpResponse(r, true)
	if err != nil {
		return fmt.Sprintf("could not dump response object: %s", err.Error())
	}
	return string(debug)
}
