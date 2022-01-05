package auth

import (
	"bytes"
	"encoding/json"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
	"github.com/pkg/errors"
)

type clientSecretJWT struct {
	tokenEndpoint string
	signer        Signer
}

func NewClientSecretJWT(tokenEndpoint string, signer Signer) Authoriser {
	return clientSecretJWT{
		tokenEndpoint: tokenEndpoint,
		signer:        signer,
	}
}

func (c clientSecretJWT) Client(response []byte) (client.Client, error) {
	var registrationResponse OBClientRegistrationResponse
	if err := json.NewDecoder(bytes.NewReader(response)).Decode(&registrationResponse); err != nil {
		return client.NewNoClient(), errors.Wrap(err, "client secret basic jwt")
	}

	return client.NewClientSecretJwt(
		registrationResponse.ClientID,
		registrationResponse.ClientSecret,
		c.tokenEndpoint,
	), nil
}

func (c clientSecretJWT) Claims() (string, error) {
	return c.signer.Claims()
}
