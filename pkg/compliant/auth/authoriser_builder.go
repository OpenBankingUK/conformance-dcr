package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
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

func (b AuthoriserBuilder) popSsas(ssas *[]string) AuthoriserBuilder {
	b.ssa = (*ssas)[0]
	if len(*ssas) > 1 {
		*ssas = (*ssas)[1:]
	} else {
		*ssas = []string{}
	}
	b.ssas = *ssas
	return b
}

// UpdateSsa - update the main ssa of the AuthoriserBuilder by popping the first one from ssas
func (b AuthoriserBuilder) UpdateSsa(ssas *[]string) AuthoriserBuilder {
	// if ssas list is empty/doesn't exist then just return not modified AuthoriserBuilder
	// if there are not enough ssas it's checked before at the early stage of dcr 32/33
	if len(*ssas) == 0 {
		return b
	} else {
		b = b.popSsas(ssas)
		return b
	}
}

// UpdateSsaAndGetSlice - UpdateSsa n times and return the generated slice of AuthoriserBuilders
func (b AuthoriserBuilder) UpdateSsaAndGetSlice(n int, ssas *[]string) []AuthoriserBuilder {
	var authoriserBuilders []AuthoriserBuilder
	for i := 0; i < n; i++ {
		newAuthoriserBuilder := b.UpdateSsa(ssas)
		authoriserBuilders = append(authoriserBuilders, newAuthoriserBuilder)
	}
	return authoriserBuilders
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
