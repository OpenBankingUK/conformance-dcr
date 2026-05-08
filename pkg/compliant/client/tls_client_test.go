package client

import (
	"io"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientTlsClientAuth(t *testing.T) {
	client := NewTlsClientAuth("id", "token")

	request, err := client.CredentialsGrantRequest()
	require.NoError(t, err)
	assert.Equal(t, "id", client.Id())
	bodyByes, err := io.ReadAll(request.Body)
	require.NoError(t, err)
	bodyDecoded, err := url.ParseQuery(string(bodyByes))
	require.NoError(t, err)

	require.Equal(t, 1, len(bodyDecoded["client_id"]))
	require.Equal(t, "id", bodyDecoded["client_id"][0])

	require.Equal(t, 1, len(bodyDecoded["grant_type"]))
	require.Equal(t, "client_credentials", bodyDecoded["grant_type"][0])

}
