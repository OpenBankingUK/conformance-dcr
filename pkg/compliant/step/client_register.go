package step

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type clientRegister struct {
	stepName        string
	client          *http.Client
	openIdConfigKey string
	responseKey     string
	ssa             string
}

func NewClientRegister(openIdConfigKey, ssa, responseKey string, httpClient *http.Client) Step {
	return clientRegister{
		stepName:        "Software client register",
		openIdConfigKey: openIdConfigKey,
		client:          httpClient,
		ssa:             ssa,
		responseKey:     responseKey,
	}
}

func (s clientRegister) Run(ctx Context) Result {
	configuration, err := ctx.GetOpenIdConfig(s.openIdConfigKey)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("getting openid config: %s", err.Error()))
	}

	response, err := s.doJwtPostRequest(configuration.RegistrationEndpoint)
	if err != nil {
		return NewFailResult(s.stepName, err.Error())
	}

	ctx.SetResponse(s.responseKey, response)

	return NewPassResult(s.stepName)
}

func (s clientRegister) doJwtPostRequest(endpoint string) (*http.Response, error) {
	body := bytes.NewBufferString(signedClaims(s.ssa))
	request, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return nil, errors.Wrap(err, "creating jwt post request")
	}
	request.Header.Add("Content-Type", "application/jwt")
	request.Header.Add("Accept", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "making jwt post request")
	}
	return response, nil
}

func signedClaims(ssa string) string {
	return ssa
}
