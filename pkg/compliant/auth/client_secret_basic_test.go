package auth

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/certs"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewClientSecretBasicAuther_Claims(t *testing.T) {
	config := openid.Configuration{}
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)
	auther := NewClientSecretBasic(config, privateKey, "ssa")

	claims, err := auther.Claims()

	assert.NoError(t, err)
	assert.NotEmpty(t, claims)
}

func TestClientSecretBasicAuther_ClientRegister_ReturnsNotImplemented(t *testing.T) {
	config := openid.Configuration{}
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)
	auther := NewClientSecretBasic(config, privateKey, "ssa")

	_, err = auther.ClientRegister([]byte{})

	assert.EqualError(t, err, "not implemented")
}
