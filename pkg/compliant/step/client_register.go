package step

import (
	"bytes"
	"fmt"
	http2 "github.com/OpenBankingUK/conformance-dcr/pkg/http"
	"github.com/pkg/errors"
	"net/http"
)

type clientRegister struct {
	stepName             string
	client               *http.Client
	registrationEndpoint string
	responseCtxKey       string
	jwtClaimsCtxKey      string
	debug                *DebugMessages
}

func NewPostClientRegister(registrationEndpoint, jwtClaimsCtxKey, responseCtxKey string, httpClient *http.Client) Step {
	return clientRegister{
		stepName:             "Software client register",
		registrationEndpoint: registrationEndpoint,
		client:               httpClient,
		jwtClaimsCtxKey:      jwtClaimsCtxKey,
		responseCtxKey:       responseCtxKey,
		debug:                NewDebug(),
	}
}

func (s clientRegister) Run(ctx Context) Result {
	s.debug.Logf("get jwt claims from ctx var: %s", s.jwtClaimsCtxKey)
	jwtClaims, err := ctx.GetString(s.jwtClaimsCtxKey)
	if err != nil {
		return s.failResult(fmt.Sprintf("getting jwt claims: %s", err.Error()))
	}

	response, err := s.doJwtPostRequest(s.registrationEndpoint, jwtClaims)
	if err != nil {
		return s.failResult(err.Error())
	}

	s.debug.Logf("setting response object in context var: %s", s.responseCtxKey)
	ctx.SetResponse(s.responseCtxKey, response)

	return NewPassResultWithDebug(s.stepName, s.debug)
}

func (s clientRegister) doJwtPostRequest(endpoint, jwtClaims string) (*http.Response, error) {
	body := bytes.NewBufferString(jwtClaims)
	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return nil, errors.Wrap(err, "creating jose post request")
	}
	req.Header.Add("Content-Type", "application/jose")
	req.Header.Add("Accept", "application/json")
	s.debug.Log(http2.DebugRequest(req))

	s.debug.Log("making request")
	response, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "making jose post request")
	}
	s.debug.Logf("request finished with response status code %d", response.StatusCode)

	return response, nil
}

func (s clientRegister) failResult(msg string) Result {
	return NewFailResultWithDebug(
		s.stepName,
		msg,
		s.debug,
	)
}
