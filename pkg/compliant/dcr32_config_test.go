package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDCR32Config(t *testing.T) {
	config := NewDCR32Config(
		openid.Configuration{},
		"ssa",
		"kid",
		[]string{""},
		&rsa.PrivateKey{},
		true,
		false,
		false,
	)

	assert.Equal(t, openid.Configuration{}, config.OpenIDConfig)
	assert.Equal(t, "ssa", config.SSA)
	assert.Equal(t, "kid", config.KID)
	assert.Equal(t, []string{""}, config.RedirectURIs)
	assert.Equal(t, &rsa.PrivateKey{}, config.PrivateKey)
	assert.True(t, config.GetImplemented)
	assert.False(t, config.PutImplemented)
	assert.False(t, config.DeleteImplemented)
}
