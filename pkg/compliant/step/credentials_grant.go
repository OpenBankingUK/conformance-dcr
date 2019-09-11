package step

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	http2 "bitbucket.org/openbankingteam/conformance-dcr/pkg/http"
	"github.com/pkg/errors"
)

type clientCredentialsGrant struct {
	client           *http.Client
	grantTokenCtxKey string
	clientCtxKey     string
	tokenEndpoint    string
	stepName         string
}

func NewClientCredentialsGrant(grantTokenCtxKey, clientCtxKey, tokenEndpoint string, httpClient *http.Client) Step {
	return clientCredentialsGrant{
		client:           httpClient,
		grantTokenCtxKey: grantTokenCtxKey,
		clientCtxKey:     clientCtxKey,
		tokenEndpoint:    tokenEndpoint,
		stepName:         fmt.Sprintf("Client credentials grant"),
	}
}

func (a clientCredentialsGrant) Run(ctx Context) Result {
	debug := NewDebug()

	softwareClient, err := ctx.GetClient(a.clientCtxKey)
	if err != nil {
		msg := fmt.Sprintf("getting software client object from context: %s", err.Error())
		return NewFailResultWithDebug(a.stepName, msg, debug)
	}
	r, err := a.requestForClient(softwareClient)
	if err != nil {
		msg := fmt.Sprintf("unable to build request object: %s", err.Error())
		return NewFailResultWithDebug(a.stepName, msg, debug)
	}

	r.Header.Set("Content-type", "application/x-www-form-urlencoded")
	debug.Log(http2.DebugRequest(r))

	response, err := a.client.Do(r)
	if err != nil {
		message := fmt.Sprintf("error making token request call: %s", err.Error())
		return NewFailResultWithDebug(a.stepName, message, debug)
	}
	debug.Log(http2.DebugResponse(response))

	if response.StatusCode != http.StatusOK {
		message := fmt.Sprintf("unexpected status code %d, should be %d", response.StatusCode, http.StatusOK)
		return NewFailResultWithDebug(a.stepName, message, debug)
	}

	var credentialsGrantResponse auth.CredentialsGrantResponse
	if err = json.NewDecoder(response.Body).Decode(&credentialsGrantResponse); err != nil {
		message := fmt.Sprintf("error decoding body content: %s", err.Error())
		return NewFailResultWithDebug(a.stepName, message, debug)
	}

	token := auth.GrantToken(credentialsGrantResponse)
	debug.Logf("setting client credentials token in context var: %s", a.grantTokenCtxKey)
	ctx.SetGrantToken(a.grantTokenCtxKey, token)

	return NewPassResultWithDebug(a.stepName, debug)
}

func credentialsGrantRequestReader() io.Reader {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "accounts openid")
	return strings.NewReader(data.Encode())
}

func (a clientCredentialsGrant) requestForClient(softwareClient client.Client) (*http.Request, error) {
	switch softwareClient.Type() {
	case client.AuthTypeClientSecretBasic:
		r, err := http.NewRequest(http.MethodPost, a.tokenEndpoint, credentialsGrantRequestReader())
		if err != nil {
			return nil, errors.Wrapf(err, "error making token request for client_secret_basic: %s", err.Error())
		}
		token, err := softwareClient.Token()
		if err != nil {
			return nil, errors.Wrapf(err, "error generating token for client_secret_basic: %s", err.Error())
		}
		r.Header.Add("Authorization", token)
		return r, nil
	case client.AuthTypePrivateKeyJwt:
		token, err := softwareClient.Token()
		if err != nil {
			return nil, errors.Wrapf(err, "error generating token for private_key_jwt: %s", err.Error())
		}
		r, err := http.NewRequest(http.MethodPost, a.tokenEndpoint, credentialGrantPrivateKeyJwtRequestReader(token))
		if err != nil {
			return nil, errors.Wrapf(err, "error making token request for private_key_jwt: %s", err.Error())
		}
		return r, nil
	}
	return nil, fmt.Errorf("auth method %s not supported", softwareClient.Type())
}

func credentialGrantPrivateKeyJwtRequestReader(token string) io.Reader {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "accounts openid")
	data.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	data.Set("client_assertion", token)
	return strings.NewReader(data.Encode())
}
