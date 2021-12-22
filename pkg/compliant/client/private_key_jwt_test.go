package client

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/url"
	"testing"
)

func TestPrivateKeyJwt(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	client := NewPrivateKeyJwt("id", "token", key, jwt.SigningMethodPS256)

	request, err := client.CredentialsGrantRequest()
	require.NoError(t, err)
	assert.Equal(t, "id", client.Id())
	bodyByes, err := ioutil.ReadAll(request.Body)
	require.NoError(t, err)
	bodyDecoded, err := url.ParseQuery(string(bodyByes))
	require.NoError(t, err)

	require.Equal(t, 1, len(bodyDecoded["client_assertion_type"]))
	require.Equal(t, "urn:ietf:params:oauth:client-assertion-type:jwt-bearer", bodyDecoded["client_assertion_type"][0])

	require.Equal(t, 1, len(bodyDecoded["grant_type"]))
	require.Equal(t, "client_credentials", bodyDecoded["grant_type"][0])
}
