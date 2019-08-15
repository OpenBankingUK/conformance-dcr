package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"encoding/json"
	"fmt"
)

type parseWellKnownRegistrationEndpoint struct {
	responseContextVar string
	openIdConfigKey    string
	stepName           string
}

func NewParseWellKnownRegistrationEndpoint(responseContextVar, openIdConfigKey string) Step {
	return parseWellKnownRegistrationEndpoint{
		responseContextVar: responseContextVar,
		openIdConfigKey:    openIdConfigKey,
		stepName:           "Decode well-known response registration endpoint",
	}
}

func (s parseWellKnownRegistrationEndpoint) Run(ctx Context) Result {
	response, err := ctx.GetResponse(s.responseContextVar)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("getting response object from context: %s", err.Error()))
	}

	defer response.Body.Close()
	config := openid.Configuration{}
	if err := json.NewDecoder(response.Body).Decode(&config); err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("reading response body: %s", err.Error()))
	}

	ctx.SetOpenIdConfig(s.openIdConfigKey, config)

	return NewPassResult(s.stepName)
}
