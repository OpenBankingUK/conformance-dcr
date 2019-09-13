package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientBasic(t *testing.T) {
	client := NewClientBasic("id", "token")

	assert.Equal(t, "id", client.Id())
	assert.Equal(t, "Basic aWQ6dG9rZW4=", client.Token())
}

func TestNoClient(t *testing.T) {
	client := NewNoClient()

	assert.Equal(t, "", client.Id())
	assert.Equal(t, "", client.Token())
}
