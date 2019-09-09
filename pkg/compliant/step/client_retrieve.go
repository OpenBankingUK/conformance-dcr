package step

import (
	"fmt"
	"net/http"
)

type clientRetrieve struct {
	client               *http.Client
	stepName             string
	clientCtxKey         string
	grantTokenCtxKey     string
	registrationEndpoint string
	responseCtxKey       string
}

func NewClientRetrieve(
	responseCtxKey, registrationEndpoint, clientCtxKey, grantTokenCtxKey string,
	httpClient *http.Client,
) Step {
	return clientRetrieve{
		stepName:             "Software client retrieve",
		client:               httpClient,
		registrationEndpoint: registrationEndpoint,
		responseCtxKey:       responseCtxKey,
		clientCtxKey:         clientCtxKey,
		grantTokenCtxKey:     grantTokenCtxKey,
	}
}

func (s clientRetrieve) Run(ctx Context) Result {
	client, err := ctx.GetClient(s.clientCtxKey)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to find client %s in context: %v", s.clientCtxKey, err))
	}

	grantToken, err := ctx.GetGrantToken(s.grantTokenCtxKey)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to find grant token %s in context: %v", s.grantTokenCtxKey, err))
	}

	endpoint := fmt.Sprintf("%s/%s", s.registrationEndpoint, client.Id)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to make request: %s", err.Error()))
	}
	req.Header.Set("Authorization", "Bearer "+grantToken.AccessToken)

	res, err := s.client.Do(req)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to call endpoint %s: %v", endpoint, err))
	}

	ctx.SetResponse(s.responseCtxKey, res)
	return NewPassResult(s.stepName)
}
