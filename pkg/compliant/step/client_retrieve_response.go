package step

import (
	"encoding/json"
	"fmt"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
)

type clientRetrieveResponse struct {
	stepName       string
	responseCtxKey string
	clientCtxKey   string
	tokenEndpoint  string
}

func NewClientRetrieveResponse(responseCtxKey, clientCtxKey, tokenEndpoint string) Step {
	return clientRetrieveResponse{
		stepName:       "Decode client retrieve response",
		responseCtxKey: responseCtxKey,
		clientCtxKey:   clientCtxKey,
		tokenEndpoint:  tokenEndpoint,
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

	ctx.SetClient(s.clientCtxKey, client.NewClientSecretBasic(
		registrationResponse.ClientID,
		registrationResponse.ClientSecret,
		s.tokenEndpoint,
	))

	return NewPassResult(s.stepName)
}
