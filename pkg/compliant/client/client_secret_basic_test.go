package client

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientBasic(t *testing.T) {
	client := NewClientSecretBasic("id", "token", "secret")

	request, err := client.CredentialsGrantRequest()
	require.NoError(t, err)
	assert.Equal(t, "id", client.Id())
	assert.Equal(t, "Basic aWQ6c2VjcmV0", request.Header.Get("Authorization"))
}
