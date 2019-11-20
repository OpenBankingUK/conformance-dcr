package auth

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	"github.com/pkg/errors"
)

type clientPrivateKeyJwt struct {
	tokenEndpoint string
	privateKey    *rsa.PrivateKey
	signer        Signer
}

func NewClientPrivateKeyJwt(tokenEndpoint string, privateKey *rsa.PrivateKey, signer Signer) Authoriser {
	return clientPrivateKeyJwt{
		tokenEndpoint: tokenEndpoint,
		privateKey:    privateKey,
		signer:        signer,
	}
}

func (c clientPrivateKeyJwt) Client(response []byte) (client.Client, error) {
	var registrationResponse OBClientRegistrationResponse
	if err := json.NewDecoder(bytes.NewReader(response)).Decode(&registrationResponse); err != nil {
		return client.NewNoClient(), errors.Wrap(err, "private key jwt client")
	}

	return client.NewClientPrivateKeyJwt(
		registrationResponse.ClientID,
		c.tokenEndpoint,
		c.privateKey,
	), nil
}

func (c clientPrivateKeyJwt) Claims() (string, error) {
	return c.signer.Claims()
}
