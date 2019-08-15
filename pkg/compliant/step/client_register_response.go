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
}

func NewClientRegisterResponse(responseCtxKey, clientCtxKey string) Step {
	return clientRegisterResponse{
		stepName:       "Decode client register response",
		responseCtxKey: responseCtxKey,
		clientCtxKey:   clientCtxKey,
	}
}

func (s clientRegisterResponse) Run(ctx Context) Result {
	response, err := ctx.GetResponse(s.responseCtxKey)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("getting response object from context: %s", err.Error()))
	}

	var registrationResponse OBClientRegistrationResponse
	if err = json.NewDecoder(response.Body).Decode(&registrationResponse); err != nil {
		return NewFailResult(s.stepName, "decoding response: "+err.Error())
	}

	ctx.SetClient(s.clientCtxKey, mapToClient(registrationResponse))

	return NewPassResult(s.stepName)
}

func mapToClient(registrationResponse OBClientRegistrationResponse) client.Client {
	return client.NewClient(
		registrationResponse.ClientID,
		registrationResponse.ClientSecret,
	)
}

type OBClientRegistrationResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
