package step

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/openid"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientRetrieve(t *testing.T) {
	// creating a stub server that expects a JWT body posted
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, fmt.Sprintf("/%s", clientID), r.URL.EscapedPath())
		_, err := w.Write([]byte(`OK`))
		require.NoError(t, err)
	}))
	defer server.Close()

	ctx := NewContext()
	ctx.SetClient("clientKey", client.NewClientSecretBasic(clientID, server.URL, clientSecret))
	ctx.SetGrantToken("grantTokenKey", auth.GrantToken{})
	step := NewClientRetrieve("responseCtxKey", server.URL, "clientKey", "grantTokenKey", server.Client())

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Software client retrieve", result.Name)
	assert.Equal(t, "", result.FailReason)

	// assert that response in now in ctx
	_, err := ctx.GetResponse("responseCtxKey")
	assert.NoError(t, err)
}

func TestNewClientRegister_HandlesMakeRequestError(t *testing.T) {
	ctx := NewContext()
	ctx.SetClient("clientKey", client.NewClientSecretBasic(clientID, "", clientSecret))
	ctx.SetGrantToken("grantTokenKey", auth.GrantToken{})
	step := NewClientRetrieve("responseCtxKey", string(rune(0x7f)), "clientKey", "grantTokenKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to make request: parse \"\\u007f/foo\": net/url: invalid control character in URL",
		result.FailReason,
	)
}

func TestNewClientRegister_HandlesExecuteRequestError(t *testing.T) {
	ctx := NewContext()
	ctx.SetClient("clientKey", client.NewClientSecretBasic(clientID, "", clientSecret))
	ctx.SetGrantToken("grantTokenKey", auth.GrantToken{})
	step := NewClientRetrieve("responseCtxKey", "localhost", "clientKey", "grantTokenKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to call endpoint localhost/foo: Get \"localhost/foo\": unsupported protocol scheme \"\"",
		result.FailReason,
	)
}

func TestNewClientRegister_HandlesErrorForClientNotFound(t *testing.T) {
	ctx := NewContext()
	registrationEndpoint := string(rune(0x7f))
	ctx.SetOpenIdConfig("openIdConfigCtxKey", openid.Configuration{
		RegistrationEndpoint: &registrationEndpoint,
		TokenEndpoint:        "",
	})
	step := NewClientRetrieve("responseCtxKey", "localhost", "clientKey", "grantTokenKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to find client clientKey in context: key not found in context",
		result.FailReason,
	)
}

func TestNewClientRegister_HandlesErrorForGrantTokenNotFound(t *testing.T) {
	ctx := NewContext()
	ctx.SetClient("clientKey", client.NewClientSecretBasic(clientID, "", clientSecret))
	step := NewClientRetrieve("responseCtxKey", "localhost", "clientKey", "grantTokenKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to find grant token grantTokenKey in context: key not found in context",
		result.FailReason,
	)
}
