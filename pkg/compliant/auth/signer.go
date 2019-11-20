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
	signingAlgorithm        string
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
	signingAlgorithm,
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
	signingMethod := jwt.SigningMethodRS256
	claims := jwt.MapClaims{
		// standard claims
		"aud": s.issuer,
		"exp": exp.Unix(),
		"jti": id.String(),
		"iat": iat.Unix(),
		"iss": s.softwareID, // TPP's unique software ID from SSA

		// metadata
		"kid":                             s.kID,
		"token_endpoint_auth_signing_alg": signingMethod.Alg(),
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
		"id_token_signed_response_alg": signingMethod.Alg(),
	}

	if s.tokenEndpointAuthMethod == "tls_client_auth" {
		if s.transportCert == nil {
			return "", errors.New("transport cert not available to get Subject")
		}
		claims["tls_client_auth_subject_dn"] = s.transportCert.Subject.ToRDNSequence().String()
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	signedJwt, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signing claims")
	}

	return signedJwt, nil
}
