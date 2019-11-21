package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
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

func Test_AuthoriserBuilder_FailsOnMissingSoftwareID(t *testing.T) {
	_, err := NewAuthoriserBuilder().
		WithSSA("ssa").
		WithKID("kid").
		Build()
	assert.EqualError(t, err, "missing softwareID from authoriser")
}

func Test_AuthoriserBuilder_FailsOnMissingPrivateKey(t *testing.T) {
	_, err := NewAuthoriserBuilder().
		WithSSA("ssa").
		WithKID("kid").
		WithSoftwareID("softwareID").
		Build()
	assert.EqualError(t, err, "missing privateKey from authoriser")
}

func Test_AuthoriserBuilder_Success(t *testing.T) {
	cert := &x509.Certificate{}

	authoriser, err := NewAuthoriserBuilder().
		WithSSA("ssa").
		WithKID("kid").
		WithSoftwareID("softwareID").
		WithPrivateKey(&rsa.PrivateKey{}).
		WithTokenEndpointAuthMethod(jwt.SigningMethodPS256).
		WithRedirectURIs([]string{"/redirect"}).
		WithTransportCert(cert).
		Build()

	assert.NoError(t, err)
	assert.Equal(t, NewAuthoriser(
		openid.Configuration{},
		"ssa",
		"kid",
		"softwareID",
		jwt.SigningMethodPS256,
		[]string{"/redirect"},
		&rsa.PrivateKey{},
		0,
		cert,
	), authoriser)
}
