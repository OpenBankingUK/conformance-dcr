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

func Test_AuthoriserBuilderSSAs_Success(t *testing.T) {
	cert := &x509.Certificate{}

	_, err := NewAuthoriserBuilder().
		WithKID("kid").
		WithIssuer("issuer").
		WithPrivateKey(&rsa.PrivateKey{}).
		WithTokenEndpointAuthMethod(jwt.SigningMethodPS256).
		WithRedirectURIs([]string{"/redirect"}).
		WithResponseTypes([]string{"code", "code id_token"}).
		WithTransportCert(cert).
		WithSSAs([]string{"ssa1", "ssa2"}).
		Build()

	assert.NoError(t, err)
}

func TestPopSSAs(t *testing.T) {
	b := AuthoriserBuilder{ssas: []string{"ssa1", "ssa2"}}
	b.popSSAs()
	assert.Equal(t, b.ssa, "ssa1")
	assert.Equal(t, b.ssas, []string{"ssa2"})
}

func TestUpdateSSASuccessForSSAsSuccess(t *testing.T) {
	b := AuthoriserBuilder{ssas: []string{"ssa1", "ssa2"}, ssasPresent: true}

	err := b.UpdateSSA()
	assert.NoError(t, err)
	assert.Equal(t, b.ssa, "ssa1")
	assert.Equal(t, b.ssas, []string{"ssa2"})

	err = b.UpdateSSA()
	assert.NoError(t, err)
	assert.Equal(t, b.ssa, "ssa2")
	assert.Equal(t, b.ssas, []string{})
}

func TestUpdateSSASuccessForSSAsFail(t *testing.T) {
	b := AuthoriserBuilder{ssas: []string{}, ssasPresent: true}

	err := b.UpdateSSA()
	assert.Error(t, err)
}

func TestUpdateSSASuccessForSSA(t *testing.T) {
	b := AuthoriserBuilder{ssa: "ssa"}

	err := b.UpdateSSA()
	assert.NoError(t, err)
	assert.Equal(t, b.ssa, "ssa")
}

func TestUpdateSSAAndGetSliceSuccess(t *testing.T) {
	b := AuthoriserBuilder{ssas: []string{"ssa1", "ssa2"}, ssasPresent: true}
	bSlice, err := b.UpdateSSAAndGetSlice(2)

	assert.NoError(t, err)
	assert.Equal(t, b.ssa, "ssa2")
	assert.Equal(t, b.ssas, []string{})

	assert.Equal(t, bSlice[0], AuthoriserBuilder{ssa: "ssa1", ssas: []string{"ssa2"}, ssasPresent: true})
	assert.Equal(t, bSlice[1], AuthoriserBuilder{ssa: "ssa2", ssas: []string{}, ssasPresent: true})
	assert.Equal(t, len(bSlice), 2)
}

func TestUpdateSSAAndGetSliceFail(t *testing.T) {
	b := AuthoriserBuilder{ssas: []string{"ssa1"}, ssasPresent: true}
	_, err := b.UpdateSSAAndGetSlice(2)

	assert.Error(t, err)
}

func TestCheckMissingSSAs(t *testing.T) {
	b := AuthoriserBuilder{missingSSAs: 0}
	err := b.CheckMissingSSAs()
	assert.NoError(t, err)

	b = AuthoriserBuilder{missingSSAs: 1}
	err = b.CheckMissingSSAs()
	assert.Error(t, err)
}
