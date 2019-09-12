package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestNewClientRegisterResponse(t *testing.T) {
	openIdConfig := openid.Configuration{TokenEndpointAuthMethodsSupported: []string{"client_secret_basic"}}
	auther := auth.NewAuthoriser(openIdConfig, "ssa", "kid", "clientId", []string{}, &rsa.PrivateKey{})

	ctx := NewContext()
	body := ioutil.NopCloser(strings.NewReader(`{"client_id": "12345", "client_secret": "54321"}`))
	ctx.SetResponse("response", &http.Response{Body: body})
	step := NewClientRegisterResponse("response", "clientCtxKey", auther)

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Decode client register response", result.Name)
	client, err := ctx.GetClient("clientCtxKey")
	require.NoError(t, err)
	r, err := client.CredentialsGrantRequest()
	require.NoError(t, err)
	assert.Equal(t, "12345", client.Id())
	assert.Equal(t, "Basic MTIzNDU6", r.Header.Get("Authorization"))
}

func TestNewClientRegisterResponse_FailsIfResponseNotFoundInContext(t *testing.T) {
	openIdConfig := openid.Configuration{TokenEndpointAuthMethodsSupported: []string{"client_secret_basic"}}
	auther := auth.NewAuthoriser(openIdConfig, "ssa", "kid", "clientId", []string{}, &rsa.PrivateKey{})
	ctx := NewContext()
	step := NewClientRegisterResponse("response", "clientCtxKey", auther)

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "getting response object from context: key not found in context", result.FailReason)
}

func TestNewClientRegisterResponse_HandlesParsingResponseObject(t *testing.T) {
	openIdConfig := openid.Configuration{TokenEndpointAuthMethodsSupported: []string{"client_secret_basic"}}
	auther := auth.NewAuthoriser(openIdConfig, "ssa", "kid", "clientId", []string{}, &rsa.PrivateKey{})
	ctx := NewContext()
	body := ioutil.NopCloser(strings.NewReader(`invalid json`))
	ctx.SetResponse("response", &http.Response{Body: body})
	step := NewClientRegisterResponse("response", "clientCtxKey", auther)

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"client register: client secret basic client: invalid character 'i' looking for beginning of value",
		result.FailReason,
	)
}
