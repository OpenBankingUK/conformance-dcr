package step

import (
	"fmt"
	"net/url"
)

type registrationEndpointValidate struct {
	stepName             string
	registrationEndpoint string
}

func NewValidateRegistrationEndpoint(registrationEndpoint string) Step {
	return registrationEndpointValidate{
		stepName:             "Registration Endpoint Validate",
		registrationEndpoint: registrationEndpoint,
	}
}

func (v registrationEndpointValidate) Run(ctx Context) Result {
	_, err := url.ParseRequestURI(v.registrationEndpoint)
	if err != nil {
		return NewFailResult(
			v.stepName,
			fmt.Sprintf("registration endpoint %s is invalid: err=%+v", v.registrationEndpoint, err),
		)
	}

	return NewPassResult(v.stepName)
}
