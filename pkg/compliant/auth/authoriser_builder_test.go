package auth

import (
	"crypto/rsa"
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
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

func Test_AuthoriserBuilder_FailsOnMissingClientID(t *testing.T) {
	_, err := NewAuthoriserBuilder().
		WithSSA("ssa").
		WithKID("kid").
		Build()
	assert.EqualError(t, err, "missing clientID from authoriser")
}

func Test_AuthoriserBuilder_FailsOnMissingPrivateKey(t *testing.T) {
	_, err := NewAuthoriserBuilder().
		WithSSA("ssa").
		WithKID("kid").
		WithClientID("clientID").
		Build()
	assert.EqualError(t, err, "missing privateKey from authoriser")
}

func Test_AuthoriserBuilder_Success(t *testing.T) {
	authoriser, err := NewAuthoriserBuilder().
		WithSSA("ssa").
		WithKID("kid").
		WithClientID("clientID").
		WithPrivateKey(&rsa.PrivateKey{}).
		Build()
	assert.Nil(t, err)
	assert.Equal(t, NewAuthoriser(
		openid.Configuration{},
		"ssa",
		"kid",
		"clientID",
		[]string{},
		&rsa.PrivateKey{},
		0,
	), authoriser)
}
