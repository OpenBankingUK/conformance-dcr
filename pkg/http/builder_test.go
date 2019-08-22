package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	rootCA, err := ioutil.ReadFile("testdata/client-sample-root-ca.pem")
	require.NoError(t, err)
	privateKey, err := ioutil.ReadFile("testdata/client-sample-key.key")
	require.NoError(t, err)
	cert, err := ioutil.ReadFile("testdata/client-sample-cert.pem")
	require.NoError(t, err)

	client, err := NewBuilder().
		WithTransportKeyPair(string(cert), string(privateKey)).
		WithRootCAs([]string{string(rootCA)}).
		Build()

	assert.NoError(t, err)
	assert.IsType(t, &http.Client{}, client)
}

func TestNewBuilder_ErrorsIfNoCertOrKey(t *testing.T) {
	rootCA, err := ioutil.ReadFile("testdata/client-sample-root-ca.pem")
	require.NoError(t, err)

	client, err := NewBuilder().
		WithRootCAs([]string{string(rootCA)}).
		Build()

	assert.EqualError(t, err, "can't build a mtls client without cert and key")
	assert.Nil(t, client)
}

func TestNewBuilder_ErrorsIfNoCertOrKeyIsInvalid(t *testing.T) {
	rootCA, err := ioutil.ReadFile("testdata/client-sample-root-ca.pem")
	require.NoError(t, err)

	client, err := NewBuilder().
		WithTransportKeyPair("", "").
		WithRootCAs([]string{string(rootCA)}).
		Build()

	assert.EqualError(
		t,
		err,
		"building mTLS http client: parse x509 key pair: tls: failed to find any PEM data in certificate input",
	)
	assert.Nil(t, client)
}

func TestNewBuilder_ErrorsIfNoRootCA(t *testing.T) {
	privateKey, err := ioutil.ReadFile("testdata/client-sample-key.key")
	require.NoError(t, err)
	cert, err := ioutil.ReadFile("testdata/client-sample-cert.pem")
	require.NoError(t, err)

	client, err := NewBuilder().
		WithTransportKeyPair(string(cert), string(privateKey)).
		Build()

	assert.EqualError(t, err, "can't build a mTLS client without rootCAs")
	assert.Nil(t, client)
}

func TestNewBuilder_ErrorsIfRootCAIsInvalid(t *testing.T) {
	privateKey, err := ioutil.ReadFile("testdata/client-sample-key.key")
	require.NoError(t, err)
	cert, err := ioutil.ReadFile("testdata/client-sample-cert.pem")
	require.NoError(t, err)

	client, err := NewBuilder().
		WithTransportKeyPair(string(cert), string(privateKey)).
		WithRootCAs([]string{"asfdasfdfs"}).
		Build()

	assert.EqualError(
		t,
		err,
		"building mTLS http client: building rootCAs certificate: 0: could not find a PEM formatted block",
	)
	assert.Nil(t, client)
}
