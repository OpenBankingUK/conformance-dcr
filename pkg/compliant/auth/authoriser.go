package auth

import (
	"crypto/rsa"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
)

// Double dispatch Signing method/Client abstract factory
type Authoriser interface {
	Claims() (string, error)
	Client(response []byte) (client.Client, error)
}

func NewAuthoriser(
	config openid.Configuration,
	ssa, kid, clientId string,
	redirectURIs []string,
	privateKey *rsa.PrivateKey,
) Authoriser {
	if sliceContains("private_key_jwt", config.TokenEndpointAuthMethodsSupported) {
		return NewClientPrivateKeyJwt(
			config.Issuer,
			config.TokenEndpoint,
			ssa,
			kid,
			clientId,
			redirectURIs,
			privateKey,
		)
	}
	if sliceContains("client_secret_basic", config.TokenEndpointAuthMethodsSupported) {
		return NewClientSecretBasic(
			config.Issuer,
			config.TokenEndpoint,
			ssa,
			kid,
			clientId,
			redirectURIs,
			privateKey,
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
