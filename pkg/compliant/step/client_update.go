package step

import (
	"bytes"
	"fmt"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
	http2 "github.com/OpenBankingUK/conformance-dcr/pkg/http"
	"github.com/pkg/errors"
	"net/http"
)

type clientUpdate struct {
	stepName             string
	client               *http.Client
	registrationEndpoint string
	responseCtxKey       string
	jwtClaimsCtxKey      string
	clientCtxKey         string
	grantTokenCtxKey     string
	debug                *DebugMessages
}

func NewClientUpdate(
	registrationEndpoint,
	jwtClaimsCtxKey,
	responseCtxKey,
	clientCtxKey,
	grantTokenCtxKey string,
	httpClient *http.Client,
) Step {
	return clientUpdate{
		stepName:             "Software client update",
		registrationEndpoint: registrationEndpoint,
		client:               httpClient,
		jwtClaimsCtxKey:      jwtClaimsCtxKey,
		responseCtxKey:       responseCtxKey,
		clientCtxKey:         clientCtxKey,
		grantTokenCtxKey:     grantTokenCtxKey,
		debug:                NewDebug(),
	}
}

func (s clientUpdate) Run(ctx Context) Result {
	s.debug.Logf("get jwt claims from ctx var: %s", s.jwtClaimsCtxKey)
	jwtClaims, err := ctx.GetString(s.jwtClaimsCtxKey)
	if err != nil {
		return s.failResult(fmt.Sprintf("getting jwt claims: %s", err.Error()))
	}

	client, err := ctx.GetClient(s.clientCtxKey)
	if err != nil {
		msg := fmt.Sprintf("unable to find client %s in context: %v", s.clientCtxKey, err)
		return NewFailResultWithDebug(s.stepName, msg, s.debug)
	}

	grantToken, err := ctx.GetGrantToken(s.grantTokenCtxKey)
	if err != nil {
		msg := fmt.Sprintf("unable to find client grant token %s in context: %v", s.grantTokenCtxKey, err)
		return NewFailResultWithDebug(s.stepName, msg, s.debug)
	}

	endpoint := fmt.Sprintf("%s/%s", s.registrationEndpoint, client.Id())
	response, err := s.doJwtPutRequest(endpoint, jwtClaims, grantToken)
	if err != nil {
		return s.failResult(err.Error())
	}

	s.debug.Logf("setting response object in context var: %s", s.responseCtxKey)
	ctx.SetResponse(s.responseCtxKey, response)

	return NewPassResultWithDebug(s.stepName, s.debug)
}

func (s clientUpdate) doJwtPutRequest(endpoint, jwtClaims string, grantToken auth.GrantToken) (*http.Response, error) {
	body := bytes.NewBufferString(jwtClaims)
	req, err := http.NewRequest(http.MethodPut, endpoint, body)
	if err != nil {
		return nil, errors.Wrap(err, "creating jose put request")
	}
	req.Header.Add("Content-Type", "application/jose")
	req.Header.Add("Accept", "application/json")

	req.Header.Set("Authorization", "Bearer "+grantToken.AccessToken)

	s.debug.Log(http2.DebugRequest(req))

	s.debug.Log("making request")
	response, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "making jose put request")
	}
	s.debug.Logf("request finished with response status code %d", response.StatusCode)

	return response, nil
}

func (s clientUpdate) failResult(msg string) Result {
	return NewFailResultWithDebug(
		s.stepName,
		msg,
		s.debug,
	)
}
