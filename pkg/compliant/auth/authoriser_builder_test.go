package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"testing"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/openid"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func Test_AuthoriserBuilder_FailsOnMissingSSA(t *testing.T) {
	_, err := NewAuthoriserBuilder().Build()
	assert.EqualError(t, err, "missing ssa from authoriser")
}

func Test_AuthoriserBuilder_FailsOnMissingKID(t *testing.T) {
	_, err := NewAuthoriserBuilder().
		WithSSA("ssa").
		Build()
	assert.EqualError(t, err, "missing kid from authoriser")
}

func Test_AuthoriserBuilder_FailsOnMissingPrivateKey(t *testing.T) {
	_, err := NewAuthoriserBuilder().
		WithSSA("ssa").
		WithKID("kid").
		WithIssuer("issuer").
		Build()
	assert.EqualError(t, err, "missing privateKey from authoriser")
}

func Test_AuthoriserBuilder_Success(t *testing.T) {
	cert := &x509.Certificate{}

	authoriser, err := NewAuthoriserBuilder().
		WithSSA("ssa").
		WithKID("kid").
		WithIssuer("issuer").
		WithPrivateKey(&rsa.PrivateKey{}).
		WithTokenEndpointAuthMethod(jwt.SigningMethodPS256).
		WithRedirectURIs([]string{"/redirect"}).
		WithResponseTypes([]string{"code", "code id_token"}).
		WithTransportCert(cert).
		Build()

	assert.NoError(t, err)
	assert.Equal(t, NewAuthoriser(
		openid.Configuration{},
		"ssa",
		"aud",
		"kid",
		"issuer",
		jwt.SigningMethodPS256,
		[]string{"/redirect"},
		[]string{"code", "code id_token"},
		&rsa.PrivateKey{},
		0,
		cert,
		"",
	), authoriser)
}
