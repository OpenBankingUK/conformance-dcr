package client

import (
	"crypto/rsa"
	"encoding/base64"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	AuthTypeClientSecretBasic = "client_secret_basic"
	AuthTypePrivateKeyJwt     = "private_key_jwt"
	AuthTypeNone              = "none"
)

type Client interface {
	Id() string
	Token() (string, error)
	Type() string
}

type clientPrivateKeyJwt struct {
	id            string
	tokenEndpoint string
	kid           string
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

func (c clientPrivateKeyJwt) Token() (string, error) {
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

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(c.privateKey)
}

func (c clientPrivateKeyJwt) Type() string {
	return AuthTypePrivateKeyJwt
}

type clientBasic struct {
	id    string
	token string
}

func NewClientBasic(id, token string) Client {
	return clientBasic{
		id:    id,
		token: token,
	}
}

func (c clientBasic) Id() string {
	return c.id
}

func (c clientBasic) Type() string {
	return AuthTypeClientSecretBasic
}

func (c clientBasic) Token() (string, error) {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(c.authClientKey())), nil
}

func (c clientBasic) authClientKey() string {
	return c.id + ":" + c.token
}

type noClient struct {
}

func NewNoClient() Client {
	return noClient{}
}

func (c noClient) Id() string {
	return ""
}

func (c noClient) Type() string {
	return AuthTypeNone
}

func (c noClient) Token() (string, error) {
	return "", nil
}
