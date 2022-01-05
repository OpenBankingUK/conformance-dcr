package step

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaims_Run(t *testing.T) {
	ctx := NewContext()
	config := openid.Configuration{TokenEndpointAuthMethodsSupported: []string{"client_secret_basic"}}
	authoriserBuilder := auth.NewAuthoriserBuilder().
		WithIssuer("issuer").
		WithKID("kid").
		WithSSA("ssa").
		WithPrivateKey(generateKey(t)).
		WithTokenEndpointAuthMethod(jwt.SigningMethodPS256).
		WithOpenIDConfig(config).
		WithJwtExpiration(time.Hour)
	step := NewClaims("jwtClaimsCtxKey", authoriserBuilder)

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
	authoriserBuilder := auth.NewAuthoriserBuilder().
		WithIssuer("softwareID").
		WithKID("kid").
		WithSSA("ssa").
		WithPrivateKey(generateKey(t)).
		WithTokenEndpointAuthMethod(jwt.SigningMethodPS256).
		WithOpenIDConfig(config).
		WithJwtExpiration(time.Hour)
	step := NewClaims("jwtClaimsCtxKey", authoriserBuilder)

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
