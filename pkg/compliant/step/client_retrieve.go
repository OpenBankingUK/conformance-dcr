package step

import (
	"fmt"
	"net/http"
)

type clientRetrieve struct {
	client             *http.Client
	stepName           string
	clientCtxKey       string
	openIDConfigCtxKey string
	responseCtxKey     string
}

func NewClientRetrieve(responseCtxKey, openIDConfigKey, clientCtxKey string, httpClient *http.Client) Step {
	return clientRetrieve{
		stepName:           "Software client retrieve",
		client:             httpClient,
		openIDConfigCtxKey: openIDConfigKey,
		responseCtxKey:     responseCtxKey,
		clientCtxKey:       clientCtxKey,
	}
}

func (s clientRetrieve) Run(ctx Context) Result {
	client, err := ctx.GetClient(s.clientCtxKey)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to find client %s in context: %v", s.clientCtxKey, err))
	}
	configuration, err := ctx.GetOpenIdConfig(s.openIDConfigCtxKey)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("getting openid config: %s", err.Error()))
	}
	endpoint := fmt.Sprintf("%s/%s", configuration.RegistrationEndpoint, client.Id)
	res, err := s.client.Get(endpoint)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to call endpoint %s: %v", endpoint, err))
	}
	ctx.SetResponse(s.responseCtxKey, res)
	return NewPassResult(s.stepName)
}
