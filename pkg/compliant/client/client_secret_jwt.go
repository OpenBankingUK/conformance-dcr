package client

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type clientSecretJwt struct {
	id            string
	tokenEndpoint string
	clientSecret  string
}

func NewClientSecretJwt(id, clientSecret, tokenEndpoint string) Client {
	return clientSecretJwt{
		id:            id,
		tokenEndpoint: tokenEndpoint,
		clientSecret:  clientSecret,
	}
}

func (c clientSecretJwt) Id() string {
	return c.id
}

func (c clientSecretJwt) CredentialsGrantRequest() (*http.Request, error) {
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

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(c.clientSecret))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to sign jwt token for client_secret_jwt")
	}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	data.Set("client_assertion", token)
	reqBody := strings.NewReader(data.Encode())
	r, err := http.NewRequest(http.MethodPost, c.tokenEndpoint, reqBody)
	if err != nil {
		return nil, errors.Wrapf(err, "error making token request for client_secret_jwt: %s", err.Error())
	}

	return r, nil
}
