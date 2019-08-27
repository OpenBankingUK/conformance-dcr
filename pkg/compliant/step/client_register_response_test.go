package step

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestNewClientRegisterResponse(t *testing.T) {
	ctx := NewContext()
	body := ioutil.NopCloser(strings.NewReader(`{"client_id": "12345", "client_secret": "54321"}`))
	ctx.SetResponse("response", &http.Response{Body: body})
	step := NewClientRegisterResponse("response", "clientCtxKey")

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Decode client register response", result.Name)
	client, err := ctx.GetClient("clientCtxKey")
	require.NoError(t, err)
	assert.Equal(t, "12345", client.Id)
	assert.Equal(t, "54321", client.Secret)
}

func TestNewClientRegisterResponse_FailsIfResponseNotFoundInContext(t *testing.T) {
	ctx := NewContext()
	step := NewClientRegisterResponse("response", "clientCtxKey")

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "getting response object from context: key not found in context", result.FailReason)
}

func TestNewClientRegisterResponse_HandlesParsingResponseObject(t *testing.T) {
	ctx := NewContext()
	body := ioutil.NopCloser(strings.NewReader(`invalid json`))
	ctx.SetResponse("response", &http.Response{Body: body})
	step := NewClientRegisterResponse("response", "clientCtxKey")

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "decoding response: invalid character 'i' looking for beginning of value", result.FailReason)
}
