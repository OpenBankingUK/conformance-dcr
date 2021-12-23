package auth

import (
	"bytes"
	"encoding/json"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
	"github.com/pkg/errors"
)

type clientSecretBasic struct {
	tokenEndpoint string
	signer        Signer
}

func NewClientSecretBasic(tokenEndpoint string, signer Signer) Authoriser {
	return clientSecretBasic{
		tokenEndpoint: tokenEndpoint,
		signer:        signer,
	}
}

func (c clientSecretBasic) Client(response []byte) (client.Client, error) {
	var registrationResponse OBClientRegistrationResponse
	if err := json.NewDecoder(bytes.NewReader(response)).Decode(&registrationResponse); err != nil {
		return client.NewNoClient(), errors.Wrap(err, "client secret basic client")
	}

	return client.NewClientSecretBasic(
		registrationResponse.ClientID,
		registrationResponse.ClientSecret,
		c.tokenEndpoint,
	), nil
}

type OBClientRegistrationResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret,omitempty"`
}

func (c clientSecretBasic) Claims() (string, error) {
	return c.signer.Claims()
}
