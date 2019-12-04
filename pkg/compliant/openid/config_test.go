package openid

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetConfig(t *testing.T) {
	body := `{
		"registration_endpoint": "http://registration_endpoint",
		"token_endpoint": "http://token_endpoint",
		"issuer": "issuer",
		"request_object_signing_alg_values_supported": ["alg1"],
		"token_endpoint_auth_methods_supported": ["alg2"]
		}`
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, err := rw.Write([]byte(body))
		require.NoError(t, err)
	}))
	defer server.Close()

	config, err := Get(server.URL, server.Client())

	assert.NoError(t, err)
	registrationEndpoint := "http://registration_endpoint"
	expected := Configuration{
		RegistrationEndpoint:              &registrationEndpoint,
		TokenEndpoint:                     "http://token_endpoint",
		ObjectSignAlgSupported:            []string{"alg1"},
		TokenEndpointAuthMethodsSupported: []string{"alg2"},
	}
	assert.Equal(t, expected, config)
}

func TestGet_HandlesNotOKStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusTeapot)
	}))
	defer server.Close()

	_, err := Get(server.URL, server.Client())

	assert.Errorf(
		t,
		err,
		"failed to GET OpenIDConfiguration config: url=%s, StatusCode=418, body=",
		server.URL,
	)
}

func TestGet_HandlesNotInvalidBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, err := rw.Write([]byte(`NOT JSON`))
		require.NoError(t, err)
	}))
	defer server.Close()

	_, err := Get(server.URL, server.Client())

	assert.EqualError(
		t,
		err,
		"invalid OpenIDConfiguration body content: invalid character 'N' looking for beginning of value",
	)
}

func TestGet_HandlesErrorGet(t *testing.T) {
	client := &http.Client{}
	_, err := Get("-", client)

	assert.EqualError(
		t,
		err,
		"Failed to GET OpenIDConfiguration: url=-: Get -: unsupported protocol scheme \"\"",
	)
}

func TestRegistrationEndpointAsString(t *testing.T) {
	endpoint := "/registration"
	c := Configuration{RegistrationEndpoint: &endpoint}

	assert.Equal(t, "/registration", c.RegistrationEndpointAsString())
}

func TestRegistrationEndpointAsString_EmptyOnNil(t *testing.T) {
	c := Configuration{}

	assert.Equal(t, "", c.RegistrationEndpointAsString())
}
