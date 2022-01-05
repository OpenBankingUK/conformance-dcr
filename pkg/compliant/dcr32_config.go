package compliant

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	http2 "net/http"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/schema"
	"github.com/OpenBankingUK/conformance-dcr/pkg/http"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/openid"
)

type DCR32Config struct {
	OpenIDConfig       openid.Configuration
	SSA                string
	KID                string
	RedirectURIs       []string
	TokenSigningMethod jwt.SigningMethod
	PrivateKey         *rsa.PrivateKey
	SecureClient       *http2.Client
	GetImplemented     bool
	PutImplemented     bool
	DeleteImplemented  bool
	AuthoriserBuilder  auth.AuthoriserBuilder
	SchemaValidator    schema.Validator
}

func NewDCR32Config(
	openIDConfig openid.Configuration,
	ssa, aud, kid, issuer string,
	redirectURIs []string,
	signingKeyPEM string,
	transportSigningKeyPEM string,
	transportCertPEM string,
	transportCertSubjectDn string,
	transportRootCAs []string,
	getImplemented bool,
	putImplemented bool,
	deleteImplemented bool,
	tlsSkipVerify bool,
	specVersion string,
) (DCR32Config, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(signingKeyPEM))
	if err != nil {
		return DCR32Config{}, errors.Wrap(err, "creating DCR32 config")
	}

	schemaValidator, err := schema.NewValidator(specVersion)
	if err != nil {
		return DCR32Config{}, errors.Wrap(err, "creating DCR32 config")
	}

	transportCert, err := certificate(transportCertPEM)
	if err != nil {
		return DCR32Config{}, errors.Wrap(err, "creating DCR32 config")
	}

	tokenSignMethod, err := responseTokenSignMethod(openIDConfig.TokenEndpointSigningAlgSupported)
	if err != nil {
		return DCR32Config{}, errors.Wrap(err, "creating DCR32 config")
	}

	responseTypes, err := responseTypeResolve(openIDConfig.ResponseTypesSupported)
	if err != nil {
		return DCR32Config{}, errors.Wrap(err, "creating DCR32 config")
	}

	// default authoriser
	authoriserBuilder := auth.NewAuthoriserBuilder().
		WithOpenIDConfig(openIDConfig).
		WithSSA(ssa).
		WithAud(aud).
		WithKID(kid).
		WithIssuer(issuer).
		WithRedirectURIs(redirectURIs).
		WithResponseTypes(responseTypes).
		WithPrivateKey(privateKey).
		WithTokenEndpointAuthMethod(tokenSignMethod).
		WithTransportCert(transportCert).
		WithTransportCertSubjectDn(transportCertSubjectDn)

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

func certificate(transportCertPEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(transportCertPEM))
	if block == nil {
		return nil, errors.New("failed making certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.New("failed making certificate")
	}
	return cert, nil
}
