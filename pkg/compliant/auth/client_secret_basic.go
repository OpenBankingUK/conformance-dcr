package auth

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type clientSecretBasic struct {
	issuer       string
	privateKey   *rsa.PrivateKey
	ssa          string
	kid          string
	clientId     string
	redirectURIs []string
}

func NewClientSecretBasic(
	issuer, ssa, kid, clientId string,
	redirectURIs []string,
	privateKey *rsa.PrivateKey,
) Authoriser {
	return clientSecretBasic{
		issuer:       issuer,
		privateKey:   privateKey,
		ssa:          ssa,
		kid:          kid,
		clientId:     clientId,
		redirectURIs: redirectURIs,
	}
}

func (c clientSecretBasic) Client(response []byte) (client.Client, error) {
	var registrationResponse OBClientRegistrationResponse
	if err := json.NewDecoder(bytes.NewReader(response)).Decode(&registrationResponse); err != nil {
		return client.NewNoClient(), errors.Wrap(err, "client secret basic client")
	}

	return NewClientBasicFromResponse(registrationResponse), nil
}

func NewClientBasicFromResponse(registrationResponse OBClientRegistrationResponse) client.Client {
	return client.NewClientBasic(
		registrationResponse.ClientID,
		registrationResponse.ClientSecret,
	)
}

type OBClientRegistrationResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret,omitempty"`
}

func (c clientSecretBasic) Claims() (string, error) {
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
			//"nbf": "",
			//"sub": "",

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
			"token_endpoint_auth_method": "client_secret_basic",
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

func (c clientSecretBasic) signClaims(token *jwt.Token) (string, error) {
	value, err := token.SignedString(c.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signing claims")
	}
	return value, nil
}
