package client

import "encoding/base64"

type Client struct {
	Id     string
	Secret string
}

func NewClient(clientId, secret string) Client {
	return Client{
		Id:     clientId,
		Secret: secret,
	}
}

func (c Client) AuthHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(c.authClientKey()))
}

func (c Client) authClientKey() string {
	return c.Id + ":" + c.Secret
}
