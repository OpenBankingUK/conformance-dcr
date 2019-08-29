package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"crypto/rand"
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClaims_Run(t *testing.T) {
	ctx := NewContext()
	config := openid.Configuration{TokenEndpointAuthMethodsSupported: []string{"client_secret_basic"}}
	step := NewClaims("jwtClaimsCtxKey", auth.NewAuthoriser(config, generateKey(t), ""))

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Generate signed software client claims", result.Name)
	assert.Equal(t, "", result.FailReason)
	claims, err := ctx.GetString("jwtClaimsCtxKey")
	assert.NoError(t, err)
	assert.NotEmpty(t, claims)
}

func TestClaims_Run_FailsOnClaimsError(t *testing.T) {
	ctx := NewContext()
	config := openid.Configuration{TokenEndpointAuthMethodsSupported: []string{""}}
	step := NewClaims("jwtClaimsCtxKey", auth.NewAuthoriser(config, nil, ""))

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "no authoriser was found for openid config", result.FailReason)
}

func generateKey(t *testing.T) *rsa.PrivateKey {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		require.NoError(t, err)
	}
	return key
}
