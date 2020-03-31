package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Signer interface {
	Claims() (string, error)
}

type jwtSigner struct {
	signingAlgorithm        jwt.SigningMethod
	ssa                     string
	issuer                  string
	audience                string
	kID                     string
	tokenEndpointAuthMethod string
	requestObjectSignAlg    string
	redirectURIs            []string
	responseTypes           *[]string
	privateKey              *rsa.PrivateKey
	jwtExpiration           time.Duration
	transportCert           *x509.Certificate
}

func NewJwtSigner(
	signingAlgorithm jwt.SigningMethod,
	ssa,
	issuer,
	audience,
	kID,
	tokenEndpointAuthMethod string,
	requestObjectSignAlg string,
	redirectURIs []string,
	responseTypes *[]string,
	privateKey *rsa.PrivateKey,
	jwtExpiration time.Duration,
	transportCert *x509.Certificate,
) Signer {
	return jwtSigner{
		signingAlgorithm:        signingAlgorithm,
		ssa:                     ssa,
		issuer:                  issuer,
		audience:                audience,
		kID:                     kID,
		tokenEndpointAuthMethod: tokenEndpointAuthMethod,
		requestObjectSignAlg:    requestObjectSignAlg,
		redirectURIs:            redirectURIs,
		responseTypes:           responseTypes,
		privateKey:              privateKey,
		jwtExpiration:           jwtExpiration,
		transportCert:           transportCert,
	}
}

func (s jwtSigner) Claims() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "generating claims")
	}

	iat := time.Now().UTC()
	exp := iat.Add(s.jwtExpiration)
	claims := jwt.MapClaims{
		// This should be the unique identifier for the ASPSP
		// issued by the issuer of the software statement.
		// An ASPSP processing the software statement may validate the
		// value of the claim and reject software statements for which the ASPSP is not the audience.
		// The value must be a Base62 encoded GUID.
		"aud": s.audience,

		"exp": exp.Unix(),
		"jti": id.String(),
		"iat": iat.Unix(),

		// Identifier for the TPP.
		// This value must be unique for each TPP registered by the issuer of the SSA.
		// The value must be a Base62 encoded GUID.
		// For SSAs issued by the OB Directory, this must be the software_id
		"iss": s.issuer,

		// metadata
		"token_endpoint_auth_signing_alg": s.signingAlgorithm.Alg(),
		"grant_types": []string{
			"authorization_code",
			"client_credentials",
		},
		"subject_type":                 "public",
		"application_type":             "web",
		"redirect_uris":                s.redirectURIs,
		"token_endpoint_auth_method":   s.tokenEndpointAuthMethod,
		"software_statement":           s.ssa,
		"scope":                        "accounts openid",
		"request_object_signing_alg":   s.requestObjectSignAlg,
		"id_token_signed_response_alg": s.signingAlgorithm.Alg(),
	}

	responseTypes, err := responseTypeResolve(s.responseTypes)
	if err != nil {
		return "", err
	}
	if responseTypes != nil {
		claims["response_types"] = responseTypes
	}

	if s.tokenEndpointAuthMethod == "tls_client_auth" {
		if s.transportCert == nil {
			return "", errors.New("transport cert not available to get Subject")
		}
		claims["tls_client_auth_subject_dn"] = s.transportCert.Subject.ToRDNSequence().String()
	}

	token := jwt.NewWithClaims(s.signingAlgorithm, claims)
	token.Header["kid"] = s.kID

	signedJwt, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signing claims")
	}

	return signedJwt, nil
}
