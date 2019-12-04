package auth

import (
	"crypto/rsa"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
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

	auther := NewAuthoriser(
		openIdConfig,
		"ssa",
		"aud",
		"kid",
		"softwareID",
		jwt.SigningMethodPS256,
		[]string{},
		&rsa.PrivateKey{},
		time.Hour,
		nil,
	)

	assert.IsType(t, clientSecretBasic{}, auther)
}

func TestNewAuther_ReturnsPrivateKeyJwt(t *testing.T) {
	openIdConfig := openid.Configuration{
		TokenEndpointAuthMethodsSupported: []string{"private_key_jwt"},
	}

	auther := NewAuthoriser(
		openIdConfig,
		"ssa",
		"aud",
		"kid",
		"softwareID",
		jwt.SigningMethodPS256,
		[]string{},
		&rsa.PrivateKey{},
		time.Hour,
		nil,
	)

	assert.IsType(t, clientPrivateKeyJwt{}, auther)
}

func TestNewAuther_ReturnsTlsClientAuth(t *testing.T) {
	openIdConfig := openid.Configuration{
		TokenEndpointAuthMethodsSupported: []string{"tls_client_auth"},
	}

	auther := NewAuthoriser(
		openIdConfig,
		"ssa",
		"aud",
		"kid",
		"softwareID",
		jwt.SigningMethodPS256,
		[]string{},
		&rsa.PrivateKey{},
		time.Hour,
		nil,
	)

	assert.IsType(t, tlsClientAuth{}, auther)
}

func TestNewAuther_ReturnsNoAuther(t *testing.T) {
	openIdConfig := openid.Configuration{
		TokenEndpointAuthMethodsSupported: []string{},
	}

	auther := NewAuthoriser(
		openIdConfig,
		"ssa",
		"aud",
		"kid",
		"softwareID",
		jwt.SigningMethodPS256,
		[]string{},
		&rsa.PrivateKey{},
		time.Hour,
		nil,
	)

	assert.IsType(t, none{}, auther)
}
