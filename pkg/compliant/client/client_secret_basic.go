package client

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

type clientSecretBasic struct {
	id            string
	tokenEndpoint string
	secret        string
}

func NewClientSecretBasic(id, tokenEndpoint, secret string) Client {
	return clientSecretBasic{
		id:            id,
		tokenEndpoint: tokenEndpoint,
		secret:        secret,
	}
}

func (c clientSecretBasic) Id() string {
	return c.id
}

func (c clientSecretBasic) CredentialsGrantRequest() (*http.Request, error) {
	token := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(c.authClientKey())))
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "openid")
	reqBody := strings.NewReader(data.Encode())
	r, err := http.NewRequest(http.MethodPost, c.tokenEndpoint, reqBody)
	if err != nil {
		return nil, errors.Wrapf(err, "error making token request for client_secret_basic: %s", err.Error())
	}
	r.Header.Add("Authorization", token)

	return r, nil
}

func (c clientSecretBasic) authClientKey() string {
	return c.id + ":" + c.secret
}
