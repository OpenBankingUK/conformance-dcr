package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"encoding/json"
	"fmt"
)

type clientRetrieveResponse struct {
	stepName       string
	responseCtxKey string
	clientCtxKey   string
}

func NewClientRetrieveResponse(responseCtxKey, clientCtxKey string) Step {
	return clientRetrieveResponse{
		stepName:       "Decode client retrieve response",
		responseCtxKey: responseCtxKey,
		clientCtxKey:   clientCtxKey,
	}
}

func (s clientRetrieveResponse) Run(ctx Context) Result {
	response, err := ctx.GetResponse(s.responseCtxKey)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("getting response object from context: %s", err.Error()))
	}

	var registrationResponse auth.OBClientRegistrationResponse
	if err = json.NewDecoder(response.Body).Decode(&registrationResponse); err != nil {
		return NewFailResult(s.stepName, "decoding response: "+err.Error())
	}

	ctx.SetClient(s.clientCtxKey, auth.NewClientBasicFromResponse(registrationResponse))

	return NewPassResult(s.stepName)
}
