package auth

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type clientSecretBasic struct {
	config     openid.Configuration
	privateKey *rsa.PrivateKey
	ssa        string
}

func NewClientSecretBasic(config openid.Configuration, privateKey *rsa.PrivateKey, ssa string) Authoriser {
	return clientSecretBasic{
		config:     config,
		privateKey: privateKey,
		ssa:        ssa,
	}
}

func (c clientSecretBasic) ClientRegister(response []byte) (client.Client, error) {
	return client.Client{}, errors.New("not implemented")
}

func (c clientSecretBasic) Claims() (string, error) {
	iat := time.Now()
	exp := iat.Add(time.Hour)
	signingMethod := jwt.SigningMethodRS256
	token := jwt.NewWithClaims(
		signingMethod,
		jwt.MapClaims{
			"kid":                             "YqL1S1MVsiknkoNpAMcXXui0VOQ",
			"token_endpoint_auth_signing_alg": signingMethod.Alg(),
			"grant_types": []string{
				"authorization_code",
				"refresh_token",
				"client_credentials",
			},
			"subject_type":     "public",
			"application_type": "web",
			"iss":              c.config.Issuer,
			"redirect_uris": []string{
				"http://redirec_url",
			},
			"token_endpoint_auth_method": "client_secret_basic",
			"aud":                        c.config.Issuer,
			"software_statement":         c.ssa,
			"scopes": []string{
				"openid",
				"accounts",
			},
			"request_object_signing_alg": "none",
			"exp":                        exp.Unix(),
			"iat":                        iat.Unix(),
			"jti":                        uuid.New().String(),
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
