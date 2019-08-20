package auth

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValueIn(t *testing.T) {
	list := []string{"one", "two", "five"}

	assert.True(t, sliceContains("one", list))
	assert.True(t, sliceContains("two", list))
	assert.True(t, sliceContains("five", list))
	assert.False(t, sliceContains("four", list))
}

func TestNewAuther_ReturnsClientSecretBasic(t *testing.T) {
	openIdConfig := openid.Configuration{
		TokenEndpointAuthMethodsSupported: []string{"client_secret_basic"},
	}

	auther := NewAuthoriser(openIdConfig, &rsa.PrivateKey{}, "ssa")

	assert.IsType(t, clientSecretBasic{}, auther)
}

func TestNewAuther_ReturnsNoAuther(t *testing.T) {
	openIdConfig := openid.Configuration{
		TokenEndpointAuthMethodsSupported: []string{},
	}

	auther := NewAuthoriser(openIdConfig, &rsa.PrivateKey{}, "ssa")

	assert.IsType(t, none{}, auther)
}
