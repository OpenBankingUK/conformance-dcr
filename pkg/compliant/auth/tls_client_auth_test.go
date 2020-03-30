package auth

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTlsClientAuth(t *testing.T) {
	a := NewTlsClientAuth("/token", mockedSigner{})

	data := []byte(`{"client_id": "12345"}`)
	client, err := a.Client(data)

	require.NoError(t, err)
	assert.Equal(t, "12345", client.Id())
}

func TestNewTlsClientAuth_ClientHandlerMarshalError(t *testing.T) {
	a := NewTlsClientAuth("/token", mockedSigner{})

	data := []byte(`{`)
	_, err := a.Client(data)

	assert.EqualError(t, err, "tls client auth: unexpected EOF")
}

func TestNewTlsClientAuth_GeneratesClaims(t *testing.T) {
	a := NewTlsClientAuth("/token", mockedSigner{})

	claims, err := a.Claims()

	require.NoError(t, err)
	assert.Equal(t, "hello", claims)
}

type mockedSigner struct{}

func (m mockedSigner) Claims() (string, error) {
	return "hello", nil
}
