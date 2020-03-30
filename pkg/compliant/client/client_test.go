package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoClient(t *testing.T) {
	client := NewNoClient()

	_, err := client.CredentialsGrantRequest()
	require.NoError(t, err)
	assert.Equal(t, "", client.Id())
}
