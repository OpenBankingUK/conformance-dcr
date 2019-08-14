package step

import (
	"encoding/json"
	"fmt"
)

type parseWellKnownRegistrationEndpoint struct {
	responseContextVar string
	storeContextVar    string
	stepName           string
}

func NewParseWellKnownRegistrationEndpoint(responseContextVar, storeContextVar string) Step {
	return parseWellKnownRegistrationEndpoint{
		responseContextVar: responseContextVar,
		storeContextVar:    storeContextVar,
		stepName:           "parse well-known response registration endpoint",
	}
}

func (s parseWellKnownRegistrationEndpoint) Run(ctx Context) Result {
	response, err := ctx.GetResponse(s.responseContextVar)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("getting response object from context: %s", err.Error()))
	}

	defer response.Body.Close()
	config := OpenIDConfiguration{}
	if err := json.NewDecoder(response.Body).Decode(&config); err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("reading response body: %s", err.Error()))
	}

	ctx.SetString(s.storeContextVar, config.RegistrationEndpoint)

	return NewPassResult(s.stepName)
}

type OpenIDConfiguration struct {
	RegistrationEndpoint string `json:"registration_endpoint"`
}
