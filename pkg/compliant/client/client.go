package client

import (
	"net/http"
)

type Client interface {
	Id() string
	CredentialsGrantRequest() (*http.Request, error)
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
