package step

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientRetrieveWithInvalidToken(t *testing.T) {
	// creating a stub server that expects a JWT body posted
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, fmt.Sprintf("/%s", clientID), r.URL.EscapedPath())
		authHeader := r.Header.Get("Authorization")
		_, err := w.Write([]byte(`OK`))
		require.NoError(t, err)
		require.Equal(t, authHeader, "Bearer foobar")
	}))
	defer server.Close()

	ctx := NewContext()
	ctx.SetClient("clientKey", client.NewClient(clientID, clientSecret))
	step := NewClientRetrieveWithInvalidToken("responseCtxKey", server.URL, "clientKey", server.Client())

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Software client retrieve with invalid token", result.Name)
	assert.Equal(t, "", result.FailReason)

	// assert that response in now in ctx
	_, err := ctx.GetResponse("responseCtxKey")
	assert.NoError(t, err)
}

func TestNewClientRetrieveWithInvalidToken_HandlesMakeRequestError(t *testing.T) {
	ctx := NewContext()
	ctx.SetClient("clientKey", client.NewClient(clientID, clientSecret))
	step := NewClientRetrieveWithInvalidToken("responseCtxKey", string(0x7f), "clientKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to make request: parse \u007f/foo: net/url: invalid control character in URL",
		result.FailReason,
	)
}

func TestNewClientRegisterWithInvalidToken_HandlesExecuteRequestError(t *testing.T) {
	ctx := NewContext()
	ctx.SetClient("clientKey", client.NewClient(clientID, clientSecret))
	step := NewClientRetrieveWithInvalidToken("responseCtxKey", "localhost", "clientKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to call endpoint localhost/foo: Get localhost/foo: unsupported protocol scheme \"\"",
		result.FailReason,
	)
}

func TestNewClientRegisterWithInvalidToken_HandlesErrorForClientNotFound(t *testing.T) {
	ctx := NewContext()
	registrationEndpoint := string(0x7f)
	ctx.SetOpenIdConfig("openIdConfigCtxKey", openid.Configuration{
		RegistrationEndpoint: &registrationEndpoint,
		TokenEndpoint:        "",
	})
	step := NewClientRetrieveWithInvalidToken("responseCtxKey", "localhost", "clientKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to find client clientKey in context: key not found in context",
		result.FailReason,
	)
}
