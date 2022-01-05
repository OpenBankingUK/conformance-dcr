package auth

import (
	"testing"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
	"github.com/stretchr/testify/assert"
)

func TestNone_Claims(t *testing.T) {
	auther := none{}

	claims, err := auther.Claims()
	assert.Error(t, err)
	assert.Equal(t, "", claims)

	c, err := auther.Client([]byte{})
	assert.Error(t, err)
	assert.Equal(t, client.NewNoClient(), c)
}
