package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientCredentialsGrant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"access_token": "takeit"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	softClient := client.NewClientSecretBasic(clientID, clientSecret, server.URL)
	ctx := NewContext()
	ctx.SetClient("clientKey", softClient)
	ctx.SetGrantToken("clientGrantKey", auth.GrantToken{})
	step := NewClientCredentialsGrant("clientGrantKey", "clientKey", server.URL, server.Client())

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Client credentials grant", result.Name)
	assert.Equal(t, "", result.FailReason)

	token, err := ctx.GetGrantToken("clientGrantKey")
	require.NoError(t, err)
	assert.Equal(t, "takeit", token.AccessToken)
}

func TestClientCredentialsGrant_HandlesClientNotFound(t *testing.T) {
	ctx := NewContext()
	step := NewClientRetrieve("responseCtxKey", "localhost", "clientKey", "grantTokenKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to find client clientKey in context: key not found in context",
		result.FailReason,
	)
}
