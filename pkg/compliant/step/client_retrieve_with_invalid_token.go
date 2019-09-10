package step

import (
	"fmt"
	"net/http"

	http2 "bitbucket.org/openbankingteam/conformance-dcr/pkg/http"
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
	debug := NewDebug()
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
	debug.Log(http2.DebugRequest(req))
	res, err := s.client.Do(req)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to call endpoint %s: %v", endpoint, err))
	}
	debug.Log(http2.DebugResponse(res))

	ctx.SetResponse(s.responseCtxKey, res)
	return NewPassResult(s.stepName)
}
