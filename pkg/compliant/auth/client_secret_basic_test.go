package auth

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/certs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewClientSecretBasicAuther_Claims(t *testing.T) {
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)
	auther := NewClientSecretBasic("issuer", "ssa", "kid", "clientId", []string{}, privateKey)

	claims, err := auther.Claims()

	assert.NoError(t, err)
	assert.NotEmpty(t, claims)
}

func TestClientSecretBasicAuther_Client_ReturnsAClient(t *testing.T) {
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)
	auther := NewClientSecretBasic("issuer", "ssa", "kid", "clientId", []string{}, privateKey)

	client, err := auther.Client([]byte(`{"client_id": "12345", "client_secret": "54321"}`))

	require.NoError(t, err)
	assert.Equal(t, "12345", client.Id())
	assert.Equal(t, "Basic MTIzNDU6NTQzMjE=", client.Token())
}
