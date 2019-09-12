package auth

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type JwtSigner struct {
	signingAlgorithm        string
	ssa                     string
	clientID                string
	issuer                  string
	kID                     string
	tokenEndpointAuthMethod string
	redirectURIs            []string
	privateKey              *rsa.PrivateKey
}

func NewJwtSigner(
	signingAlgorithm,
	ssa,
	clientID,
	issuer,
	kID,
	tokenEndpointAuthMethod string,
	redirectURIs []string,
	privateKey *rsa.PrivateKey,
) JwtSigner {
	return JwtSigner{
		signingAlgorithm:        signingAlgorithm,
		ssa:                     ssa,
		clientID:                clientID,
		issuer:                  issuer,
		kID:                     kID,
		tokenEndpointAuthMethod: tokenEndpointAuthMethod,
		redirectURIs:            redirectURIs,
		privateKey:              privateKey,
	}
}

func (s JwtSigner) Claims() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "generating claims")
	}

	iat := time.Now().UTC()
	exp := iat.Add(time.Hour)
	signingMethod := jwt.SigningMethodRS256
	token := jwt.NewWithClaims(
		signingMethod,
		jwt.MapClaims{
			// standard claims
			"aud": s.issuer,
			"exp": exp.Unix(),
			"jti": id.String(),
			"iat": iat.Unix(),
			"iss": s.clientID,

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
		},
	)

	signedJwt, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signing claims")
	}

	return signedJwt, nil
}
