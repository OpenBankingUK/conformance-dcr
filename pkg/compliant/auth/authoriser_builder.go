package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"github.com/dgrijalva/jwt-go"
)

type AuthoriserBuilder struct {
	config                  openid.Configuration
	ssa, aud, kID, issuer   string
	tokenEndpointAuthMethod jwt.SigningMethod
	redirectURIs            []string
	responseTypes           *[]string
	privateKey              *rsa.PrivateKey
	jwtExpiration           time.Duration
	transportCert           *x509.Certificate
}

func NewAuthoriserBuilder() AuthoriserBuilder {
	return AuthoriserBuilder{
		tokenEndpointAuthMethod: jwt.SigningMethodPS256,
		jwtExpiration:           time.Hour,
	}
}

func (b AuthoriserBuilder) WithTransportCert(transportCert *x509.Certificate) AuthoriserBuilder {
	b.transportCert = transportCert
	return b
}

func (b AuthoriserBuilder) WithOpenIDConfig(cfg openid.Configuration) AuthoriserBuilder {
	b.config = cfg
	return b
}

func (b AuthoriserBuilder) WithSSA(ssa string) AuthoriserBuilder {
	b.ssa = ssa
	return b
}

func (b AuthoriserBuilder) WithIssuer(issuer string) AuthoriserBuilder {
	b.issuer = issuer
	return b
}

func (b AuthoriserBuilder) WithAud(aud string) AuthoriserBuilder {
	b.aud = aud
	return b
}

func (b AuthoriserBuilder) WithKID(kID string) AuthoriserBuilder {
	b.kID = kID
	return b
}

func (b AuthoriserBuilder) WithTokenEndpointAuthMethod(alg jwt.SigningMethod) AuthoriserBuilder {
	b.tokenEndpointAuthMethod = alg
	return b
}

func (b AuthoriserBuilder) WithRedirectURIs(redirectURIs []string) AuthoriserBuilder {
	b.redirectURIs = redirectURIs
	return b
}

func (b AuthoriserBuilder) WithResponseTypes(responseTypes *[]string) AuthoriserBuilder {
	b.responseTypes = responseTypes
	return b
}

func (b AuthoriserBuilder) WithPrivateKey(privateKey *rsa.PrivateKey) AuthoriserBuilder {
	b.privateKey = privateKey
	return b
}

func (b AuthoriserBuilder) WithJwtExpiration(jwtExpiration time.Duration) AuthoriserBuilder {
	b.jwtExpiration = jwtExpiration
	return b
}

func (b AuthoriserBuilder) Build() (Authoriser, error) {
	if b.ssa == "" {
		return none{}, errors.New("missing ssa from authoriser")
	}
	if b.kID == "" {
		return none{}, errors.New("missing kid from authoriser")
	}
	if b.privateKey == nil {
		return none{}, errors.New("missing privateKey from authoriser")
	}
	return NewAuthoriser(
		b.config,
		b.ssa,
		b.aud,
		b.kID,
		b.issuer,
		b.tokenEndpointAuthMethod,
		b.redirectURIs,
		b.responseTypes,
		b.privateKey,
		b.jwtExpiration,
		b.transportCert,
	), nil
}
