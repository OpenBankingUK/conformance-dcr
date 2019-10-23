package auth

import (
	"crypto/rsa"
	"errors"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"github.com/dgrijalva/jwt-go"
)

type AuthoriserBuilder struct {
	config                                        openid.Configuration
	ssa, kID, softwareID, tokenEndpointAuthMethod string
	redirectURIs                                  []string
	privateKey                                    *rsa.PrivateKey
	jwtExpiration                                 time.Duration
}

func NewAuthoriserBuilder() AuthoriserBuilder {
	return AuthoriserBuilder{tokenEndpointAuthMethod: jwt.SigningMethodPS256.Alg()}
}

func (b AuthoriserBuilder) WithOpenIDConfig(cfg openid.Configuration) AuthoriserBuilder {
	b.config = cfg
	return b
}

func (b AuthoriserBuilder) WithSSA(ssa string) AuthoriserBuilder {
	b.ssa = ssa
	return b
}

func (b AuthoriserBuilder) WithSoftwareID(softwareID string) AuthoriserBuilder {
	b.softwareID = softwareID
	return b
}

func (b AuthoriserBuilder) WithKID(kID string) AuthoriserBuilder {
	b.kID = kID
	return b
}

func (b AuthoriserBuilder) WithTokenEndpointAuthMethod(alg string) AuthoriserBuilder {
	b.tokenEndpointAuthMethod = alg
	return b
}

func (b AuthoriserBuilder) WithRedirectURIs(redirectURIs []string) AuthoriserBuilder {
	b.redirectURIs = redirectURIs
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
	if b.softwareID == "" {
		return none{}, errors.New("missing softwareID from authoriser")
	}
	if b.privateKey == nil {
		return none{}, errors.New("missing privateKey from authoriser")
	}
	return NewAuthoriser(
		b.config,
		b.ssa,
		b.kID,
		b.softwareID,
		b.tokenEndpointAuthMethod,
		b.redirectURIs,
		b.privateKey,
		b.jwtExpiration,
	), nil
}
