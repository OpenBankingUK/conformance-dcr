package schema

import (
	"bytes"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestResponseValidator32_ValidateInvalidPayload(t *testing.T) {
	validator := responseValidator32{}
	reader := bytes.NewReader([]byte(`{`))

	failures := validator.Validate(reader)

	assert.Len(t, failures, 1)
}

func TestResponseValidator32_ValidateEmpty(t *testing.T) {
	validator := responseValidator32{}
	data := []byte(`{}`)
	reader := bytes.NewReader(data)

	failures := validator.Validate(reader)

	assert.Len(t, failures, 9)
}

func TestResponseValidator32_ValidateResponse(t *testing.T) {
	validator := responseValidator32{}
	reader, err := os.Open("testdata/response.json")
	require.NoError(t, err)

	failures := validator.Validate(reader)

	assert.Len(t, failures, 0)
}

func TestResponseValidationSlice(t *testing.T) {
	data := struct {
		urls []string
	}{
		urls: []string{"aaa", "bbb"},
	}

	err := validation.Validate(
		data.urls,
		validation.Required,
		validation.Each(validation.Length(1, 256), isOBURLValidationRule()),
	)

	require.Error(t, err)
}

func TestResponseValidationTokenAuthMethodPrivateKey(t *testing.T) {
	method := "private_key_jwt"
	client := OBClientRegistrationResponseSchema32{
		TokenEndpointAuthMethod: &method,
	}

	failures := validateTokenEndpointMethod32(client)

	expectedFailures := []Failure{
		"token_endpoint_auth_signing_alg MUST be specified if " +
			"token_endpoint_auth_method is private_key_jwt or client_secret_jwt",
	}
	assert.Equal(t, expectedFailures, failures)
}

func TestResponseValidationTokenAuthMethodClientSecret(t *testing.T) {
	method := "client_secret_jwt"
	client := OBClientRegistrationResponseSchema32{
		TokenEndpointAuthMethod: &method,
	}

	failures := validateTokenEndpointMethod32(client)

	expectedFailures := []Failure{
		"token_endpoint_auth_signing_alg MUST be specified if " +
			"token_endpoint_auth_method is private_key_jwt or client_secret_jwt",
	}
	assert.Equal(t, expectedFailures, failures)
}

func TestResponseValidationTokenAuthMethodTLSClientAuth(t *testing.T) {
	method := "tls_client_auth"
	client := OBClientRegistrationResponseSchema32{
		TokenEndpointAuthMethod: &method,
	}

	failures := validateTokenEndpointMethod32(client)

	expectedFailures := []Failure{
		"tls_client_auth_subject_dn MUST be set if " +
			"token_endpoint_auth_method is set to tls_client_auth",
	}
	assert.Equal(t, expectedFailures, failures)
}
