package client

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Client interface {
	Id() string
	CredentialsGrantRequest() (*http.Request, error)
}

type clientPrivateKeyJwt struct {
	id            string
	tokenEndpoint string
	privateKey    *rsa.PrivateKey
}

func NewClientPrivateKeyJwt(id, tokenEndpoint string, privateKey *rsa.PrivateKey) Client {
	return clientPrivateKeyJwt{
		id:            id,
		tokenEndpoint: tokenEndpoint,
		privateKey:    privateKey,
	}
}

func (c clientPrivateKeyJwt) Id() string {
	return c.id
}

func (c clientPrivateKeyJwt) CredentialsGrantRequest() (*http.Request, error) {
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

type clientBasic struct {
	id            string
	tokenEndpoint string
	secret        string
}

func NewClientBasic(id, tokenEndpoint, secret string) Client {
	return clientBasic{
		id:            id,
		tokenEndpoint: tokenEndpoint,
		secret:        secret,
	}
}

func (c clientBasic) Id() string {
	return c.id
}

func (c clientBasic) CredentialsGrantRequest() (*http.Request, error) {
	token := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(c.authClientKey())))
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	reqBody := strings.NewReader(data.Encode())
	r, err := http.NewRequest(http.MethodPost, c.tokenEndpoint, reqBody)
	if err != nil {
		return nil, errors.Wrapf(err, "error making token request for client_secret_basic: %s", err.Error())
	}
	r.Header.Add("Authorization", token)

	return r, nil
}

func (c clientBasic) authClientKey() string {
	return c.id + ":" + c.secret
}

type noClient struct {
}

func NewNoClient() Client {
	return noClient{}
}

func (c noClient) Id() string {
	return ""
}

func (c noClient) CredentialsGrantRequest() (*http.Request, error) {
	return nil, nil
}
