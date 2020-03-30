package step

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	http2 "bitbucket.org/openbankingteam/conformance-dcr/pkg/http"
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
	r, err := softwareClient.CredentialsGrantRequest()
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
