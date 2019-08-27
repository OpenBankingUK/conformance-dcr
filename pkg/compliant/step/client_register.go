package step

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httputil"
)

type clientRegister struct {
	stepName           string
	client             *http.Client
	openIdConfigCtxKey string
	responseCtxKey     string
	jwtClaimsCtxKey    string
	debug              *DebugMessages
}

func NewPostClientRegister(openIdConfigCtxKey, jwtClaimsCtxKey, responseCtxKey string, httpClient *http.Client) Step {
	return clientRegister{
		stepName:           "Software client register",
		openIdConfigCtxKey: openIdConfigCtxKey,
		client:             httpClient,
		jwtClaimsCtxKey:    jwtClaimsCtxKey,
		responseCtxKey:     responseCtxKey,
		debug:              NewDebug(),
	}
}

func (s clientRegister) Run(ctx Context) Result {
	s.debug.Logf("get openid config from ctx var: %s", s.openIdConfigCtxKey)
	configuration, err := ctx.GetOpenIdConfig(s.openIdConfigCtxKey)
	if err != nil {
		return s.failResult(fmt.Sprintf("getting openid config: %s", err.Error()))
	}

	s.debug.Logf("get jwt claims from ctx var: %s", s.jwtClaimsCtxKey)
	jwtClaims, err := ctx.GetString(s.jwtClaimsCtxKey)
	if err != nil {
		return s.failResult(fmt.Sprintf("getting jwt claims: %s", err.Error()))
	}

	response, err := s.doJwtPostRequest(configuration.RegistrationEndpoint, jwtClaims)
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
		return nil, errors.Wrap(err, "creating jwt post request")
	}
	req.Header.Add("Content-Type", "application/jwt")
	req.Header.Add("Accept", "application/json")
	s.debugRequest(req)

	s.debug.Log("making request")
	response, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "making jwt post request")
	}
	s.debug.Logf("request finished with response status code %d", response.StatusCode)

	return response, nil
}

func (s clientRegister) debugRequest(req *http.Request) {
	debug, err := httputil.DumpRequest(req, true)
	if err != nil {
		s.debug.Logf("cant debug request object: %s", err.Error())
	} else {
		s.debug.Logf("request built: %s", string(debug))
	}
}

func (s clientRegister) failResult(msg string) Result {
	return NewFailResultWithDebug(
		s.stepName,
		msg,
		s.debug,
	)
}
