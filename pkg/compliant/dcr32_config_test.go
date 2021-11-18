package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestNewDCR32Config(t *testing.T) {
	privateKeyPEM, err := ioutil.ReadFile("testdata/client-sample-key.key")
	require.NoError(t, err)

	certPEM, err := ioutil.ReadFile("testdata/client-sample-cert.pem")
	require.NoError(t, err)

	certRootPEM, err := ioutil.ReadFile("testdata/client-sample-root-ca.pem")
	require.NoError(t, err)

	config, err := NewDCR32Config(
		openid.Configuration{},
		"ssa",
		"aud",
		"kid",
		"ssaId",
		[]string{"/redirect"},
		string(privateKeyPEM),
		string(privateKeyPEM),
		string(certPEM),
		"",
		[]string{string(certRootPEM)},
		true,
		false,
		false,
		false,
		"3.2",
	)
	require.NoError(t, err)

	assert.Equal(t, openid.Configuration{}, config.OpenIDConfig)
	assert.Equal(t, "ssa", config.SSA)
	assert.Equal(t, "kid", config.KID)
	assert.Equal(t, []string{"/redirect"}, config.RedirectURIs)
	assert.True(t, config.GetImplemented)
	assert.False(t, config.PutImplemented)
	assert.False(t, config.DeleteImplemented)
}
