package auth

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type clientPrivateKeyJwt struct {
	issuer        string
	tokenEndpoint string
	privateKey    *rsa.PrivateKey
	ssa           string
	kid           string
	clientId      string
	redirectURIs  []string
}

func NewClientPrivateKeyJwt(
	issuer, tokenEndpoint, ssa, kid, clientId string,
	redirectURIs []string,
	privateKey *rsa.PrivateKey,
) Authoriser {
	return clientPrivateKeyJwt{
		issuer:        issuer,
		tokenEndpoint: tokenEndpoint,
		privateKey:    privateKey,
		ssa:           ssa,
		kid:           kid,
		clientId:      clientId,
		redirectURIs:  redirectURIs,
	}
}

func (c clientPrivateKeyJwt) Client(response []byte) (client.Client, error) {
	var registrationResponse OBClientRegistrationResponse
	if err := json.NewDecoder(bytes.NewReader(response)).Decode(&registrationResponse); err != nil {
		return client.NewNoClient(), errors.Wrap(err, "private key jwt client")
	}

	return client.NewClientPrivateKeyJwt(
		registrationResponse.ClientID,
		c.tokenEndpoint,
		c.privateKey,
	), nil
}

func (c clientPrivateKeyJwt) Claims() (string, error) {
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
			"aud": c.issuer,
			"exp": exp.Unix(),
			"jti": id.String(),
			"iat": iat.Unix(),
			"iss": c.clientId,

			// metadata
			"kid":                             c.kid,
			"token_endpoint_auth_signing_alg": signingMethod.Alg(),
			"grant_types": []string{
				"authorization_code",
				"client_credentials",
			},
			"subject_type":               "public",
			"application_type":           "web",
			"redirect_uris":              c.redirectURIs,
			"token_endpoint_auth_method": "private_key_jwt",
			"software_statement":         c.ssa,
			"scope":                      "accounts openid",
			"request_object_signing_alg": "none",
			"response_types": []string{
				"code",
				"code id_token",
			},
			"id_token_signed_response_alg": signingMethod.Alg(),
		},
	)
	return c.signClaims(token)
}

func (c clientPrivateKeyJwt) signClaims(token *jwt.Token) (string, error) {
	value, err := token.SignedString(c.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signing claims")
	}
	return value, nil
}
