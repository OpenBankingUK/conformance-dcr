package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"github.com/dgrijalva/jwt-go"
	"time"

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
	ssa, kid, softwareID string, tokenEndpointAuthMethod jwt.SigningMethod,
	redirectURIs []string,
	privateKey *rsa.PrivateKey,
	jwtExpiration time.Duration,
	transportCert *x509.Certificate,
) Authoriser {
	if sliceContains("tls_client_auth", config.TokenEndpointAuthMethodsSupported) {
		return NewTlsClientAuth(
			config.TokenEndpoint,
			NewJwtSigner(
				tokenEndpointAuthMethod,
				ssa,
				softwareID,
				config.Issuer,
				kid,
				"tls_client_auth",
				redirectURIs,
				privateKey,
				jwtExpiration,
				transportCert,
			),
		)
	}
	if sliceContains("private_key_jwt", config.TokenEndpointAuthMethodsSupported) {
		return NewClientPrivateKeyJwt(
			config.TokenEndpoint,
			privateKey,
			NewJwtSigner(
				tokenEndpointAuthMethod,
				ssa,
				softwareID,
				config.Issuer,
				kid,
				"private_key_jwt",
				redirectURIs,
				privateKey,
				jwtExpiration,
				transportCert,
			),
		)
	}
	if sliceContains("client_secret_jwt", config.TokenEndpointAuthMethodsSupported) {
		return NewClientSecretJWT(
			config.TokenEndpoint,
			NewJwtSigner(
				tokenEndpointAuthMethod,
				ssa,
				softwareID,
				config.Issuer,
				kid,
				"client_secret_jwt",
				redirectURIs,
				privateKey,
				jwtExpiration,
				transportCert,
			),
		)
	}
	if sliceContains("client_secret_basic", config.TokenEndpointAuthMethodsSupported) {
		return NewClientSecretBasic(
			config.TokenEndpoint,
			NewJwtSigner(
				tokenEndpointAuthMethod,
				ssa,
				softwareID,
				config.Issuer,
				kid,
				"client_secret_basic",
				redirectURIs,
				privateKey,
				jwtExpiration,
				transportCert,
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
