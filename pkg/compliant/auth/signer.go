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
	softwareID              string
	issuer                  string
	kID                     string
	tokenEndpointAuthMethod string
	redirectURIs            []string
	privateKey              *rsa.PrivateKey
	jwtExpiration           time.Duration
	transportCert           *x509.Certificate
}

func NewJwtSigner(
	signingAlgorithm jwt.SigningMethod,
	ssa,
	softwareID,
	issuer,
	kID,
	tokenEndpointAuthMethod string,
	redirectURIs []string,
	privateKey *rsa.PrivateKey,
	jwtExpiration time.Duration,
	transportCert *x509.Certificate,
) Signer {
	return jwtSigner{
		signingAlgorithm:        signingAlgorithm,
		ssa:                     ssa,
		softwareID:              softwareID,
		issuer:                  issuer,
		kID:                     kID,
		tokenEndpointAuthMethod: tokenEndpointAuthMethod,
		redirectURIs:            redirectURIs,
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
		// standard claims
		"aud": s.issuer,
		"exp": exp.Unix(),
		"jti": id.String(),
		"iat": iat.Unix(),

		// Identifier for the TPP.
		// This value must be unique for each TPP registered by the issuer of the SSA.
		// The value must be a Base62 encoded GUID.
		// For SSAs issued by the OB Directory, this must be the software_id
		"iss": s.softwareID,

		// metadata
		"token_endpoint_auth_signing_alg": s.signingAlgorithm.Alg(),
		"grant_types": []string{
			"authorization_code",
			"client_credentials",
		},
		"subject_type":               "public",
		"application_type":           "web",
		"redirect_uris":              s.redirectURIs,
		"token_endpoint_auth_method": s.tokenEndpointAuthMethod,
		"software_statement":         s.ssa,
		"scope":                      "accounts openid",
		"request_object_signing_alg": "none",
		"response_types": []string{
			"code",
			"code id_token",
		},
		"id_token_signed_response_alg": s.signingAlgorithm.Alg(),
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
