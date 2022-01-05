package auth

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	"errors"
)

type none struct{}

func (c none) Claims() (string, error) {
	return "", errors.New("no authoriser was found for openid config")
}

func (c none) Client(response []byte) (client.Client, error) {
	return client.NewNoClient(), errors.New("no authoriser was found for openid config")
}
