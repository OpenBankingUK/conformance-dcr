package http

import (
	"crypto/tls"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewMATLSClient(t *testing.T) {
	// Bootstrap the tests with required keys certificates
	rootCAs, err := rootCASFromFile("testdata/client-sample-root-ca.pem")
	if err != nil {
		t.Fatalf("load root CAs from file: %s", err.Error())
	}
	rootCAPool := rootCAPoolFromCerts(rootCAs)

	clientCerts, err := clientCertsFromFile(
		"testdata/client-sample-key.key", "testdata/client-sample-cert.pem")
	if err != nil {
		t.Fatalf("Create client certs: %s", err.Error())
	}

	config := MATLSConfig{
		InsecureSkipVerify: false,
		ClientCerts:        clientCerts,
		RootCAs:            rootCAs,
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
	wantErr := false

	got, err := NewMATLSClient(config)
	if (err != nil) != wantErr {
		t.Errorf("NewMATLSClient() error = %v, wantErr %v", err, wantErr)
		return
	}

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
