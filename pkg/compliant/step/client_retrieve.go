package step

import (
	http2 "bitbucket.org/openbankingteam/conformance-dcr/pkg/http"
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
	debug := NewDebug()

	client, err := ctx.GetClient(s.clientCtxKey)
	if err != nil {
		msg := fmt.Sprintf("unable to find client %s in context: %v", s.clientCtxKey, err)
		return NewFailResultWithDebug(s.stepName, msg, debug)
	}

	grantToken, err := ctx.GetGrantToken(s.grantTokenCtxKey)
	if err != nil {
		msg := fmt.Sprintf("unable to find grant token %s in context: %v", s.grantTokenCtxKey, err)
		return NewFailResultWithDebug(s.stepName, msg, debug)
	}

	endpoint := fmt.Sprintf("%s/%s", s.registrationEndpoint, client.Id())
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		msg := fmt.Sprintf("unable to make request: %s", err.Error())
		return NewFailResultWithDebug(s.stepName, msg, debug)
	}
	req.Header.Set("Authorization", "Bearer "+grantToken.AccessToken)

	debug.Log(http2.DebugRequest(req))
	res, err := s.client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("unable to call endpoint %s: %v", endpoint, err)
		return NewFailResultWithDebug(s.stepName, msg, debug)
	}

	ctx.SetResponse(s.responseCtxKey, res)
	return NewPassResultWithDebug(s.stepName, debug)
}
