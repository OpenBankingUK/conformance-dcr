package client

import (
	"encoding/base64"
)

type Client interface {
	Id() string
	Token() string
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

func (c clientBasic) Token() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(c.authClientKey()))
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

func (c noClient) Token() string {
	return ""
}
