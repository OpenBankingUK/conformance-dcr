package auth

import (
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewClientPrivateKeyJwt(t *testing.T) {
	privateKey := &rsa.PrivateKey{}
	a := NewClientPrivateKeyJwt("/token", privateKey, mockedSigner{})

	data := []byte(`{"client_id": "12345"}`)
	client, err := a.Client(data)

	require.NoError(t, err)
	assert.Equal(t, "12345", client.Id())
}

func TestNewClientPrivateKeyJwt_ClientHandlerMarshalError(t *testing.T) {
	privateKey := &rsa.PrivateKey{}
	a := NewClientPrivateKeyJwt("/token", privateKey, mockedSigner{})

	data := []byte(`{`)
	_, err := a.Client(data)

	assert.EqualError(t, err, "private key jwt client: unexpected EOF")
}

func TestNewClientPrivateKeyJwt_GeneratesClaims(t *testing.T) {
	privateKey := &rsa.PrivateKey{}
	a := NewClientPrivateKeyJwt("/token", privateKey, mockedSigner{})

	claims, err := a.Claims()

	require.NoError(t, err)
	assert.Equal(t, "hello", claims)
}
