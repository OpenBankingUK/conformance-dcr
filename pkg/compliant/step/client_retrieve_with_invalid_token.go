package step

import (
	"fmt"
	"net/http"
)

type clientRetrieveWithInvalidToken struct {
	client               *http.Client
	stepName             string
	clientCtxKey         string
	registrationEndpoint string
	responseCtxKey       string
}

func NewClientRetrieveWithInvalidToken(
	responseCtxKey, registrationEndpoint, clientCtxKey string,
	httpClient *http.Client,
) Step {
	return clientRetrieveWithInvalidToken{
		stepName:             "Software client retrieve with invalid token",
		client:               httpClient,
		registrationEndpoint: registrationEndpoint,
		responseCtxKey:       responseCtxKey,
		clientCtxKey:         clientCtxKey,
	}
}

func (s clientRetrieveWithInvalidToken) Run(ctx Context) Result {
	client, err := ctx.GetClient(s.clientCtxKey)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to find client %s in context: %v", s.clientCtxKey, err))
	}

	endpoint := fmt.Sprintf("%s/%s", s.registrationEndpoint, client.Id)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to make request: %s", err.Error()))
	}
	req.Header.Set("Authorization", "Bearer foobar")

	res, err := s.client.Do(req)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to call endpoint %s: %v", endpoint, err))
	}

	ctx.SetResponse(s.responseCtxKey, res)
	return NewPassResult(s.stepName)
}
