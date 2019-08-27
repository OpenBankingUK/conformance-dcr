package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	"encoding/json"
	"fmt"
)

type clientRegisterResponse struct {
	stepName       string
	responseCtxKey string
	clientCtxKey   string
	debug          *DebugMessages
}

func NewClientRegisterResponse(responseCtxKey, clientCtxKey string) Step {
	return clientRegisterResponse{
		stepName:       "Decode client register response",
		responseCtxKey: responseCtxKey,
		clientCtxKey:   clientCtxKey,
		debug:          NewDebug(),
	}
}

func (s clientRegisterResponse) Run(ctx Context) Result {
	s.debug.Logf("get response object from ctx var: %s", s.responseCtxKey)
	response, err := ctx.GetResponse(s.responseCtxKey)
	if err != nil {
		return s.failResult(fmt.Sprintf("getting response object from context: %s", err.Error()))
	}

	s.debug.Log("decoding response body to OBClientRegistrationResponse object")
	var registrationResponse OBClientRegistrationResponse
	if err = json.NewDecoder(response.Body).Decode(&registrationResponse); err != nil {
		return s.failResult("decoding response: " + err.Error())
	}

	s.debug.Logf("setting software client in context var: %s", s.clientCtxKey)
	ctx.SetClient(s.clientCtxKey, mapToClient(registrationResponse))

	return NewPassResultWithDebug(s.stepName, s.debug)
}

func (s clientRegisterResponse) failResult(msg string) Result {
	return NewFailResultWithDebug(
		s.stepName,
		msg,
		s.debug,
	)
}

func mapToClient(registrationResponse OBClientRegistrationResponse) client.Client {
	return client.NewClient(
		registrationResponse.ClientID,
		registrationResponse.ClientSecret,
	)
}

type OBClientRegistrationResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret,omitempty"`
}
