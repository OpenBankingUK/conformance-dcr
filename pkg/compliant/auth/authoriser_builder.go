package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/openid"
	"github.com/dgrijalva/jwt-go"
)

type AuthoriserBuilder struct {
	config                  openid.Configuration
	ssa, aud, kID, issuer   string
	tokenEndpointSignMethod jwt.SigningMethod
	redirectURIs            []string
	responseTypes           []string
	privateKey              *rsa.PrivateKey
	jwtExpiration           time.Duration
	transportCert           *x509.Certificate
	transportCertSubjectDn  string
	ssas                    []string
	ssasPresent             bool
	missingSSAs             int
}

func NewAuthoriserBuilder() AuthoriserBuilder {
	return AuthoriserBuilder{
		jwtExpiration: time.Hour,
	}
}

func (b AuthoriserBuilder) WithTransportCert(transportCert *x509.Certificate) AuthoriserBuilder {
	b.transportCert = transportCert
	return b
}

func (b AuthoriserBuilder) WithTransportCertSubjectDn(transportSubjectDn string) AuthoriserBuilder {
	b.transportCertSubjectDn = transportSubjectDn
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

func (b AuthoriserBuilder) WithSSAs(ssas []string) AuthoriserBuilder {
	b.ssas = ssas
	b = b.checkSSAsPresent()
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
	b.tokenEndpointSignMethod = alg
	return b
}

func (b AuthoriserBuilder) WithRedirectURIs(redirectURIs []string) AuthoriserBuilder {
	b.redirectURIs = redirectURIs
	return b
}

func (b AuthoriserBuilder) WithResponseTypes(responseTypes []string) AuthoriserBuilder {
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

func (b AuthoriserBuilder) WithSSAsPresent(ssasPresent bool) AuthoriserBuilder {
	b.ssasPresent = ssasPresent
	return b
}

func (b AuthoriserBuilder) checkSSAsPresent() AuthoriserBuilder {
	if len(b.ssas) > 0 {
		b.ssasPresent = true
	}
	return b
}

func (b *AuthoriserBuilder) popSSAs() {
	b.ssa = (b.ssas)[0]
	if len(b.ssas) > 1 {
		b.ssas = (b.ssas)[1:]
	} else {
		b.ssas = []string{}
	}
}

// UpdateSSA - update the main ssa of the AuthoriserBuilder by popping the first one from ssas
func (b *AuthoriserBuilder) UpdateSSA() error {
	if !b.ssasPresent {
		return nil
	}

	if len(b.ssas) == 0 {
		b.missingSSAs += 1
		return errors.New("not enough SSAs")
	}

	if len(b.ssas) > 0 {
		b.popSSAs()
	}

	return nil
}

// UpdateSSAAndGetSlice - UpdateSsa n times and return the generated slice of AuthoriserBuilders
func (b *AuthoriserBuilder) UpdateSSAAndGetSlice(n int) ([]AuthoriserBuilder, error) {
	var authoriserBuilders []AuthoriserBuilder

	for i := 0; i < n; i++ {
		err := b.UpdateSSA()
		if err != nil {
			b.missingSSAs += n - i - 1
			return nil, err
		}

		authoriserBuilders = append(authoriserBuilders, *b)
	}
	return authoriserBuilders, nil
}

// CheckMissingSSAs - Check if b.missingSSAs was updated (default = 0)
func (b AuthoriserBuilder) CheckMissingSSAs() error {
	if b.missingSSAs > 0 {
		return errors.New(fmt.Sprintf("invalid amount of SSAs provided in the config - missing: %d", b.missingSSAs))
	}
	return nil
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
	if b.tokenEndpointSignMethod == nil {
		return none{}, errors.New("missing token endpoint signing method from authoriser")
	}
	return NewAuthoriser(
		b.config,
		b.ssa,
		b.aud,
		b.kID,
		b.issuer,
		b.tokenEndpointSignMethod,
		b.redirectURIs,
		b.responseTypes,
		b.privateKey,
		b.jwtExpiration,
		b.transportCert,
		b.transportCertSubjectDn,
	), nil
}
