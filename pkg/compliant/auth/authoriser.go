package auth

import (
	"crypto/rsa"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"github.com/dgrijalva/jwt-go"
)

// Double dispatch Signing method/Client abstract factory
type Authoriser interface {
	Claims() (string, error)
	Client(response []byte) (client.Client, error)
}

func NewAuthoriser(
	config openid.Configuration,
	ssa, kid, softwareID string,
	redirectURIs []string,
	privateKey *rsa.PrivateKey,
	jwtExpiration time.Duration,
) Authoriser {
	if sliceContains("private_key_jwt", config.TokenEndpointAuthMethodsSupported) {
		return NewClientPrivateKeyJwt(
			config.Issuer,
			config.TokenEndpoint,
			ssa,
			kid,
			redirectURIs,
			privateKey,
			NewJwtSigner(
				jwt.SigningMethodPS256.Alg(),
				ssa,
				softwareID,
				config.Issuer,
				kid,
				"private_key_jwt",
				redirectURIs,
				privateKey,
				jwtExpiration,
			),
		)
	}
	if sliceContains("client_secret_basic", config.TokenEndpointAuthMethodsSupported) {
		return NewClientSecretBasic(
			config.Issuer,
			config.TokenEndpoint,
			ssa,
			kid,
			redirectURIs,
			privateKey,
			NewJwtSigner(
				jwt.SigningMethodPS256.Alg(),
				ssa,
				softwareID,
				config.Issuer,
				kid,
				"client_secret_basic",
				redirectURIs,
				privateKey,
				jwtExpiration,
			),
		)
	}
	return none{}
}

func sliceContains(value string, list []string) bool {
	for _, item := range list {
		if value == item {
			return true
		}
	}
	return false
}
