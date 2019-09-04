package step

import (
	"fmt"
	"net/http"
)

type clientDelete struct {
	client               *http.Client
	stepName             string
	clientCtxKey         string
	registrationEndpoint string
	responseCtxKey       string
}

func NewClientDelete(responseCtxKey, registrationEndpoint, clientCtxKey string, httpClient *http.Client) Step {
	return clientDelete{
		stepName:             "Software client delete",
		client:               httpClient,
		registrationEndpoint: registrationEndpoint,
		responseCtxKey:       responseCtxKey,
		clientCtxKey:         clientCtxKey,
	}
}

func (s clientDelete) Run(ctx Context) Result {
	client, err := ctx.GetClient(s.clientCtxKey)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to find client %s in context: %v", s.clientCtxKey, err))
	}

	url := fmt.Sprintf("%s/%s", s.registrationEndpoint, client.Id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to create request %s: %v", url, err))
	}

	res, err := s.client.Do(req)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to call endpoint %s: %v", url, err))
	}

	ctx.SetResponse(s.responseCtxKey, res)
	return NewPassResult(s.stepName)
}
