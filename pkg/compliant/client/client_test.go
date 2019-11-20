package client

import (
	"crypto/rand"
	"crypto/rsa"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientBasic(t *testing.T) {
	client := NewClientBasic("id", "token", "secret")

	request, err := client.CredentialsGrantRequest()
	require.NoError(t, err)
	assert.Equal(t, "id", client.Id())
	assert.Equal(t, "Basic aWQ6c2VjcmV0", request.Header.Get("Authorization"))
}

func TestClientPrivateKeyJwt(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	client := NewClientPrivateKeyJwt("id", "token", key)

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

func TestClientTlsClientAuth(t *testing.T) {
	client := NewTlsClientAuth("id", "token")

	request, err := client.CredentialsGrantRequest()
	require.NoError(t, err)
	assert.Equal(t, "id", client.Id())
	bodyByes, err := ioutil.ReadAll(request.Body)
	require.NoError(t, err)
	bodyDecoded, err := url.ParseQuery(string(bodyByes))
	require.NoError(t, err)

	require.Equal(t, 1, len(bodyDecoded["client_id"]))
	require.Equal(t, "id", bodyDecoded["client_id"][0])

	require.Equal(t, 1, len(bodyDecoded["grant_type"]))
	require.Equal(t, "client_credentials", bodyDecoded["grant_type"][0])
}

func TestNoClient(t *testing.T) {
	client := NewNoClient()

	_, err := client.CredentialsGrantRequest()
	require.NoError(t, err)
	assert.Equal(t, "", client.Id())
}
