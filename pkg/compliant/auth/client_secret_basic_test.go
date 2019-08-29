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

func TestClientSecretBasicAuther_ClientRegister_ReturnsNotImplemented(t *testing.T) {
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)
	auther := NewClientSecretBasic("issuer", "ssa", "kid", "clientId", []string{}, privateKey)

	_, err = auther.ClientRegister([]byte{})

	assert.EqualError(t, err, "not implemented")
}
