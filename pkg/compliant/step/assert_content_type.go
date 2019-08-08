package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"
	"fmt"
)

type assertContentType struct {
	responseContextVar string
	contentType        string
	order              int
	stepName           string
}

func NewAssertContentType(order int, responseContextVar string, contentType string) Step {
	return assertContentType{
		responseContextVar: responseContextVar,
		contentType:        contentType,
		order:              order,
		stepName:           fmt.Sprintf("Content-Type header is %s", contentType),
	}
}

func (a assertContentType) Run(ctx context.Context) Result {
	response, err := ctx.GetResponse(a.responseContextVar)
	if err != nil {
		return NewFailResult(a.stepName, fmt.Sprintf("getting response object from context: %s", err.Error()))
	}

	_, ok := response.Header["Content-Type"]
	if !ok {
		return NewFailResult(a.stepName, "Content-Type header is not present")
	}

	contentType := response.Header.Get("Content-Type")
	if contentType != a.contentType {
		return NewFailResult(a.stepName, fmt.Sprintf("Content-Type is '%s'", contentType))
	}

	return NewPassResult(a.stepName)
}

func (a assertContentType) Order() int {
	return a.order
}
