package step

import (
	"fmt"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/schema"
	http2 "github.com/OpenBankingUK/conformance-dcr/pkg/http"
)

type clientRetrieveSchema struct {
	stepName       string
	responseCtxKey string
	validator      schema.Validator
}

func NewClientRetrieveSchema(responseCtxKey string, validator schema.Validator) Step {
	return clientRetrieveSchema{
		stepName:       "Validate client response schema",
		responseCtxKey: responseCtxKey,
		validator:      validator,
	}
}

func (s clientRetrieveSchema) Run(ctx Context) Result {
	debug := NewDebug()

	debug.Logf("get response object from ctx var: %s", s.responseCtxKey)
	response, err := ctx.GetResponse(s.responseCtxKey)
	if err != nil {
		return NewFailResultWithDebug(s.stepName, fmt.Sprintf("getting response object from context: %s", err.Error()), debug)
	}

	debug.Log(http2.DebugResponse(response))

	debug.Log("cloning response body object")
	body, bodyCopy, err := http2.DrainBody(response.Body)
	if err != nil {
		return NewFailResultWithDebug(s.stepName, fmt.Sprintf("copy body from response: %s", err.Error()), debug)
	}
	response.Body = body

	failures := s.validator.Validate(bodyCopy)
	if len(failures) > 0 {
		msg := string(failures[0])
		for _, failure := range failures {
			msg = fmt.Sprintf("%s, %s", msg, failure)
		}
		return NewFailResultWithDebug(s.stepName, "schema invalid: "+msg, debug)
	}

	err = bodyCopy.Close()
	if err != nil {
		return NewFailResultWithDebug(s.stepName, fmt.Sprintf("closing body clone: %s", err.Error()), debug)
	}

	return NewPassResultWithDebug(s.stepName, debug)
}
