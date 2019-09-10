package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	http2 "bitbucket.org/openbankingteam/conformance-dcr/pkg/http"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

	r, err := http.NewRequest(http.MethodPost, a.tokenEndpoint, credentialsGrantRequestReader())
	if err != nil {
		debug.Log(http2.DebugRequest(r))
		message := fmt.Sprintf("error making token request: %s", err.Error())
		return NewFailResultWithDebug(a.stepName, message, debug)
	}
	r.Header.Add("Authorization", softwareClient.Token())
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
