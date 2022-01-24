package auth

import (
	"bytes"
	"encoding/json"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
	"github.com/pkg/errors"
)

type tlsClientAuth struct {
	tokenEndpoint string
	signer        Signer
}

func NewTlsClientAuth(tokenEndpoint string, singer Signer) Authoriser {
	return tlsClientAuth{
		tokenEndpoint: tokenEndpoint,
		signer:        singer,
	}
}

func (c tlsClientAuth) Client(response []byte) (client.Client, error) {
	var registrationResponse OBClientRegistrationResponse
	if err := json.NewDecoder(bytes.NewReader(response)).Decode(&registrationResponse); err != nil {
		return client.NewNoClient(), errors.Wrap(err, "tls client auth")
	}

	return client.NewTlsClientAuth(
		registrationResponse.ClientID,
		c.tokenEndpoint,
	), nil
}

func (c tlsClientAuth) Claims() (string, error) {
	return c.signer.Claims()
}
