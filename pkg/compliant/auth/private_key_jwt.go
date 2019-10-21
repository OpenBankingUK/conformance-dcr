package auth

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	"github.com/pkg/errors"
)

type clientPrivateKeyJwt struct {
	issuer        string
	tokenEndpoint string
	privateKey    *rsa.PrivateKey
	ssa           string
	kid           string
	redirectURIs  []string
	jwtSigner     JwtSigner
}

func NewClientPrivateKeyJwt(
	issuer, tokenEndpoint, ssa, kid string,
	redirectURIs []string,
	privateKey *rsa.PrivateKey,
	jwtSigner JwtSigner,
) Authoriser {
	return clientPrivateKeyJwt{
		issuer:        issuer,
		tokenEndpoint: tokenEndpoint,
		privateKey:    privateKey,
		ssa:           ssa,
		kid:           kid,
		redirectURIs:  redirectURIs,
		jwtSigner:     jwtSigner,
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
	return c.jwtSigner.Claims()
}
