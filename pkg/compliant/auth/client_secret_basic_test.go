package auth

import (
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/certs"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientSecretBasicAuther_Claims(t *testing.T) {
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)
	auther := NewClientSecretBasic(
		"tokenEndpoint",
		NewJwtSigner(
			jwt.SigningMethodRS256,
			"ssa",
			"issuer",
			"aud",
			"kid",
			"private_key_jwt",
			"none",
			[]string{},
			privateKey,
			time.Hour,
			nil,
		),
	)

	claims, err := auther.Claims()

	assert.NoError(t, err)
	assert.NotEmpty(t, claims)
}

func TestClientSecretBasicAuther_Client_ReturnsAClient(t *testing.T) {
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)
	auther := NewClientSecretBasic(
		"tokenEndpoint",
		NewJwtSigner(
			jwt.SigningMethodRS256,
			"ssa",
			"issuer",
			"aud",
			"kid",
			"private_key_jwt",
			"none",
			[]string{},
			privateKey,
			time.Hour,
			nil,
		),
	)

	client, err := auther.Client([]byte(`{"client_id": "12345", "client_secret": "54321"}`))
	require.NoError(t, err)
	r, err := client.CredentialsGrantRequest()
	require.NoError(t, err)
	assert.Equal(t, "12345", client.Id())
	assert.Equal(t, "Basic MTIzNDU6dG9rZW5FbmRwb2ludA==", r.Header.Get("Authorization"))
}
