package compliant

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponseTokenSignMethod(t *testing.T) {
	method, err := responseTokenSignMethod(nil)
	assert.NoError(t, err)
	assert.Equal(t, method, jwt.SigningMethodPS256)

	method, err = responseTokenSignMethod(&[]string{"PS256"})
	assert.NoError(t, err)
	assert.Equal(t, method, jwt.SigningMethodPS256)

	method, err = responseTokenSignMethod(&[]string{"RS256", "PS256"})
	assert.NoError(t, err)
	assert.Equal(t, method, jwt.SigningMethodPS256)

	method, err = responseTokenSignMethod(&[]string{"RS256"})
	assert.EqualError(t, err, "PS256 token sign method not found")
	assert.Nil(t, method)
}
