package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"github.com/dgrijalva/jwt-go"
	"time"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/openid"
)

// Double dispatch Signing method/Client abstract factory
type Authoriser interface {
	Claims() (string, error)
	Client(response []byte) (client.Client, error)
}

func NewAuthoriser(
	config openid.Configuration,
	ssa, aud, kid, issuer string, tokenEndpointSignMethod jwt.SigningMethod,
	redirectURIs []string,
	responseTypes []string,
	privateKey *rsa.PrivateKey,
	jwtExpiration time.Duration,
	transportCert *x509.Certificate,
	transportSubjectDn string,
) Authoriser {
	requestObjectSignAlg := "none"
	if len(config.RequestObjectSignAlgSupported) > 0 {
		requestObjectSignAlg = config.RequestObjectSignAlgSupported[0]
	}

	if sliceContains("tls_client_auth", config.TokenEndpointAuthMethodsSupported) {
		return NewTlsClientAuth(
			config.TokenEndpoint,
			NewJwtSigner(
				tokenEndpointSignMethod,
				ssa,
				issuer,
				aud,
				kid,
				"tls_client_auth",
				requestObjectSignAlg,
				redirectURIs,
				responseTypes,
				privateKey,
				jwtExpiration,
				transportCert,
				transportSubjectDn,
			),
		)
	}
	if sliceContains("private_key_jwt", config.TokenEndpointAuthMethodsSupported) {
		return NewClientPrivateKeyJwt(
			config.TokenEndpoint,
			tokenEndpointSignMethod,
			privateKey,
			NewJwtSigner(
				tokenEndpointSignMethod,
				ssa,
				issuer,
				aud,
				kid,
				"private_key_jwt",
				requestObjectSignAlg,
				redirectURIs,
				responseTypes,
				privateKey,
				jwtExpiration,
				transportCert,
				transportSubjectDn,
			),
		)
	}
	if sliceContains("client_secret_jwt", config.TokenEndpointAuthMethodsSupported) {
		return NewClientSecretJWT(
			config.TokenEndpoint,
			NewJwtSigner(
				tokenEndpointSignMethod,
				ssa,
				issuer,
				aud,
				kid,
				"client_secret_jwt",
				requestObjectSignAlg,
				redirectURIs,
				responseTypes,
				privateKey,
				jwtExpiration,
				transportCert,
				transportSubjectDn,
			),
		)
	}
	if sliceContains("client_secret_basic", config.TokenEndpointAuthMethodsSupported) {
		return NewClientSecretBasic(
			config.TokenEndpoint,
			NewJwtSigner(
				tokenEndpointSignMethod,
				ssa,
				issuer,
				aud,
				kid,
				"client_secret_basic",
				requestObjectSignAlg,
				redirectURIs,
				responseTypes,
				privateKey,
				jwtExpiration,
				transportCert,
				transportSubjectDn,
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
