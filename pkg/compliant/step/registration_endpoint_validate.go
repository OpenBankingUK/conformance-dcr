package step

import (
	"fmt"
	"net/url"
)

type registrationEndpointValidate struct {
	stepName             string
	registrationEndpoint *string
}

func NewValidateRegistrationEndpoint(registrationEndpoint *string) Step {
	return registrationEndpointValidate{
		stepName:             "Registration Endpoint Validate",
		registrationEndpoint: registrationEndpoint,
	}
}

func (v registrationEndpointValidate) Run(ctx Context) Result {
	if v.registrationEndpoint == nil {
		return NewFailResult(
			v.stepName,
			fmt.Sprintf("registration endpoint is missing"),
		)
	}
	_, err := url.ParseRequestURI(*v.registrationEndpoint)
	if err != nil {
		debug := NewDebug()
		debug.Logf("invalid URL err=%+v", err)

		return NewFailResultWithDebug(
			v.stepName,
			fmt.Sprintf("registration endpoint %s is invalid", *v.registrationEndpoint),
			debug,
		)
	}

	return NewPassResult(v.stepName)
}
