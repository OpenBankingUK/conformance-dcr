package step

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestNewParseWellKnownRegistrationEndpoint(t *testing.T) {
	ctx := NewContext()
	body := ioutil.NopCloser(strings.NewReader(`{"registration_endpoint": "hal"}`))
	ctx.SetResponse("response", &http.Response{Body: body})
	step := NewParseWellKnownRegistrationEndpoint("response", "registration_endpoint")

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "parse well-known response registration endpoint", result.Name)
	r, err := ctx.GetString("registration_endpoint")
	require.NoError(t, err)
	assert.Equal(t, "hal", r)
}

func TestNewParseWellKnownRegistrationEndpoint_FailsIfResponseNotFoundInContext(t *testing.T) {
	ctx := NewContext()
	step := NewParseWellKnownRegistrationEndpoint("response", "registration_endpoint")

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "getting response object from context: key not found in context", result.Message)
}

func TestNewParseWellKnownRegistrationEndpoint_HandlesParsingResponseObject(t *testing.T) {
	ctx := NewContext()
	body := ioutil.NopCloser(strings.NewReader(`invalid json`))
	ctx.SetResponse("response", &http.Response{Body: body})
	step := NewParseWellKnownRegistrationEndpoint("response", "registration_endpoint")

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "reading response body: invalid character 'i' looking for beginning of value", result.Message)
}
