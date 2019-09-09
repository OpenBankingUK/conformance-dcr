package step

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
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
