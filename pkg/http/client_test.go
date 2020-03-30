package http

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRootCASCertificate(t *testing.T) {
	cert, err := RootCASCertificate([]byte(""))

	assert.EqualError(t, err, "could not find a PEM formatted block")
	assert.Nil(t, cert)
}

func TestTlsClientCertFromFile_HandlesKeyFileError(t *testing.T) {
	certs, err := TlsCertFromFile(
		"wrongfile",
		"testdata/client-sample-cert.pem",
	)

	assert.EqualError(t, err, "tls cert from file: open wrongfile: no such file or directory")
	assert.Nil(t, certs)
}

func TestTlsClientCertFromFile_HandlesCertFileError(t *testing.T) {
	certs, err := TlsCertFromFile(
		"testdata/client-sample-key.key",
		"wrongfile",
	)

	assert.EqualError(t, err, "tls cert from file: open wrongfile: no such file or directory")
	assert.Nil(t, certs)
}

func TestRootCASFromFile_HandlesFileError(t *testing.T) {
	cert, err := RootCASFromFile("wrongfile")

	assert.EqualError(t, err, "rootCAs from file: open wrongfile: no such file or directory")
	assert.Nil(t, cert)
}

func TestNewMATLSClient(t *testing.T) {
	// Bootstrap the tests with required keys certificates
	rootCAs, err := RootCASFromFile("testdata/client-sample-root-ca.pem")
	require.NoError(t, err)

	rootCAPool := RootCAPoolFromCerts([]*x509.Certificate{rootCAs})

	clientCerts, err := TlsCertFromFile(
		"testdata/client-sample-key.key",
		"testdata/client-sample-cert.pem",
	)
	require.NoError(t, err)

	config := MATLSConfig{
		InsecureSkipVerify: false,
		ClientCerts:        clientCerts,
		RootCAs:            []*x509.Certificate{rootCAs},
		TLSMinVersion:      tls.VersionTLS12,
	}
	wantClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				Certificates:       clientCerts,
				MinVersion:         tls.VersionTLS12,
				Renegotiation:      tls.RenegotiateFreelyAsClient,
				RootCAs:            rootCAPool,
			},
		},
	}

	got, err := NewMATLSClient(config)
	require.NoError(t, err)

	trsActual, ok := got.Transport.(*http.Transport)
	assert.True(t, ok)
	trsExpected, ok := wantClient.Transport.(*http.Transport)
	assert.True(t, ok)

	assert.Equal(t, trsExpected.TLSClientConfig.MinVersion, trsActual.TLSClientConfig.MinVersion)
	assert.Equal(t, trsExpected.TLSClientConfig.InsecureSkipVerify, trsActual.TLSClientConfig.InsecureSkipVerify)
	assert.Equal(t, trsExpected.TLSClientConfig.Renegotiation, trsActual.TLSClientConfig.Renegotiation)
	assert.Equal(t, trsExpected.TLSClientConfig.Certificates, trsActual.TLSClientConfig.Certificates)

	expectedRootCAs := trsExpected.TLSClientConfig.RootCAs.Subjects()
	actualRootCAs := trsActual.TLSClientConfig.RootCAs.Subjects()
	assert.Equal(t, expectedRootCAs, actualRootCAs)
}
