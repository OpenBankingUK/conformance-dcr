package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/schema"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/http"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	http2 "net/http"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
)

type DCR32Config struct {
	OpenIDConfig      openid.Configuration
	SSA               string
	KID               string
	RedirectURIs      []string
	PrivateKey        *rsa.PrivateKey
	SecureClient      *http2.Client
	GetImplemented    bool
	PutImplemented    bool
	DeleteImplemented bool
	AuthoriserBuilder auth.AuthoriserBuilder
	SchemaValidator   schema.Validator
}

func NewDCR32Config(
	openIDConfig openid.Configuration,
	ssa, aud, kid, issuer string,
	redirectURIs []string,
	signingKeyPEM string,
	transportSigningKeyPEM string,
	transportCertPEM string,
	transportRootCAs []string,
	getImplemented bool,
	putImplemented bool,
	deleteImplemented bool,
	tlsSkipVerify bool,
) (DCR32Config, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(signingKeyPEM))
	if err != nil {
		return DCR32Config{}, errors.Wrap(err, "creating DCR32 config")
	}

	const responseSchemaVersion = "3.2"
	schemaValidator, err := schema.NewValidator(responseSchemaVersion)
	if err != nil {
		return DCR32Config{}, errors.Wrap(err, "creating DCR32 config")
	}

	block, _ := pem.Decode([]byte(transportCertPEM))
	if block == nil {
		return DCR32Config{}, errors.New("failed to parse certificate PEM")
	}
	transportCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return DCR32Config{}, errors.Wrap(err, "creating DCR32 config")
	}

	authoriserBuilder := auth.NewAuthoriserBuilder().
		WithOpenIDConfig(openIDConfig).
		WithSSA(ssa).
		WithAud(aud).
		WithKID(kid).
		WithIssuer(issuer).
		WithRedirectURIs(redirectURIs).
		WithResponseTypes(openIDConfig.ResponseTypesSupported).
		WithPrivateKey(privateKey).
		WithTransportCert(transportCert)

	secureClient, err := http.NewBuilder().
		WithRootCAs(transportRootCAs).
		WithTransportKeyPair(transportCertPEM, transportSigningKeyPEM).
		WithTlsSkipVerify(tlsSkipVerify).
		Build()
	if err != nil {
		return DCR32Config{}, errors.Wrap(err, "creating DCR32 config")
	}

	return DCR32Config{
		OpenIDConfig:      openIDConfig,
		SSA:               ssa,
		KID:               kid,
		RedirectURIs:      redirectURIs,
		PrivateKey:        privateKey,
		SecureClient:      secureClient,
		GetImplemented:    getImplemented,
		PutImplemented:    putImplemented,
		DeleteImplemented: deleteImplemented,
		AuthoriserBuilder: authoriserBuilder,
		SchemaValidator:   schemaValidator,
	}, nil
}
