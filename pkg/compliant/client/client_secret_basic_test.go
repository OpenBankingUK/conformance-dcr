package client

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/url"
	"testing"
)

func TestClientBasic(t *testing.T) {
	client := NewClientSecretBasic("id", "token", "secret")

	request, err := client.CredentialsGrantRequest()
	require.NoError(t, err)
	assert.Equal(t, "id", client.Id())
	assert.Equal(t, "Basic aWQ6c2VjcmV0", request.Header.Get("Authorization"))

	bodyByes, err := ioutil.ReadAll(request.Body)
	require.NoError(t, err)

	bodyDecoded, err := url.ParseQuery(string(bodyByes))
	require.NoError(t, err)

	require.Equal(t, 1, len(bodyDecoded["grant_type"]))
	require.Equal(t, "client_credentials", bodyDecoded["grant_type"][0])

	require.Equal(t, 1, len(bodyDecoded["scope"]))
	require.Equal(t, "openid", bodyDecoded["scope"][0])
}
