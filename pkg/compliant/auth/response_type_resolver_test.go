package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponseTypeResolver(t *testing.T) {
	responseTypes, err := responseTypeResolve(nil)
	assert.NoError(t, err)
	assert.Nil(t, responseTypes)

	responseTypes, err = responseTypeResolve(&[]string{"code"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"code"}, responseTypes)

	responseTypes, err = responseTypeResolve(&[]string{"code id_token"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"code id_token"}, responseTypes)

	responseTypes, err = responseTypeResolve(&[]string{"code id_token", "code"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"code id_token", "code"}, responseTypes)

	responseTypes, err = responseTypeResolve(&[]string{"code", "code id_token", "id_token"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"code", "code id_token"}, responseTypes)

	responseTypes, err = responseTypeResolve(&[]string{"id_token"})
	assert.EqualError(t, err, "supported response types must contain `code` and/or `code id_token`")
	assert.Nil(t, responseTypes)
}
