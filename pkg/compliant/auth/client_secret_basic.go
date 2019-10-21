package auth

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	"github.com/pkg/errors"
)

type clientSecretBasic struct {
	issuer        string
	tokenEndpoint string
	privateKey    *rsa.PrivateKey
	ssa           string
	kid           string
	redirectURIs  []string
	jwtSigner     JwtSigner
}

func NewClientSecretBasic(
	issuer, tokenEndpoint, ssa, kid string,
	redirectURIs []string,
	privateKey *rsa.PrivateKey,
	jwtSigner JwtSigner,
) Authoriser {
	return clientSecretBasic{
		issuer:        issuer,
		tokenEndpoint: tokenEndpoint,
		privateKey:    privateKey,
		ssa:           ssa,
		kid:           kid,
		redirectURIs:  redirectURIs,
		jwtSigner:     jwtSigner,
	}
}

func (c clientSecretBasic) Client(response []byte) (client.Client, error) {
	var registrationResponse OBClientRegistrationResponse
	if err := json.NewDecoder(bytes.NewReader(response)).Decode(&registrationResponse); err != nil {
		return client.NewNoClient(), errors.Wrap(err, "client secret basic client")
	}

	return client.NewClientBasic(
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
	return c.jwtSigner.Claims()
}
