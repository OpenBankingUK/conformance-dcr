package client

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type privateKeyJwt struct {
	id            string
	tokenEndpoint string
	privateKey    *rsa.PrivateKey
}

func NewPrivateKeyJwt(id, tokenEndpoint string, privateKey *rsa.PrivateKey) Client {
	return privateKeyJwt{
		id:            id,
		tokenEndpoint: tokenEndpoint,
		privateKey:    privateKey,
	}
}

func (c privateKeyJwt) Id() string {
	return c.id
}

func (c privateKeyJwt) CredentialsGrantRequest() (*http.Request, error) {
	now := time.Now()
	iat := now.Unix()
	exp := now.Add(30 * time.Minute).Unix()
	jti := uuid.New().String()
	claims := jwt.MapClaims{
		"iss": c.id,
		"sub": c.id,
		"aud": c.tokenEndpoint,
		"iat": iat,
		"exp": exp,
		"jti": jti,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(c.privateKey)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to sign jwt token for private_key_jwt")
	}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	data.Set("client_assertion", token)
	reqBody := strings.NewReader(data.Encode())
	r, err := http.NewRequest(http.MethodPost, c.tokenEndpoint, reqBody)
	if err != nil {
		return nil, errors.Wrapf(err, "error making token request for private_key_jwt: %s", err.Error())
	}

	return r, nil
}
