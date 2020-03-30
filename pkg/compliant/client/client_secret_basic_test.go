package client

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/url"
	"testing"
)

func TestClientBasic(t *testing.T) {
	client := NewClientSecretBasic("id", "secret", "http://endpoint")

	expectedTokenHeader := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("id:secret")))

	request, err := client.CredentialsGrantRequest()
	require.NoError(t, err)
	assert.Equal(t, "id", client.Id())
	assert.Equal(t, expectedTokenHeader, request.Header.Get("Authorization"))

	bodyByes, err := ioutil.ReadAll(request.Body)
	require.NoError(t, err)

	bodyDecoded, err := url.ParseQuery(string(bodyByes))
	require.NoError(t, err)

	require.Equal(t, 1, len(bodyDecoded["grant_type"]))
	require.Equal(t, "client_credentials", bodyDecoded["grant_type"][0])

	require.Equal(t, 1, len(bodyDecoded["scope"]))
	require.Equal(t, "openid", bodyDecoded["scope"][0])
}
