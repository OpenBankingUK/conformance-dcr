package http

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func clientCertsFromFile(keyPath, certPath string) ([]tls.Certificate, error) {
	keyBlock, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, errors.Wrap(err, "read key file")
	}
	certBlock, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, errors.Wrap(err, "read cert file")
	}

	crt, err := tls.X509KeyPair(certBlock, keyBlock)
	if err != nil {
		return nil, errors.Wrap(err, "parse x509 key pair")
	}
	return []tls.Certificate{crt}, nil
}

func rootCASFromFile(path string) ([]*x509.Certificate, error) {
	pemBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	// Currently only support one block
	var block *pem.Block
	block, _ = pem.Decode(pemBytes)

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parse x509 certificate")
	}

	return []*x509.Certificate{cert}, nil
}

func rootCAPoolFromCerts(certs []*x509.Certificate) *x509.CertPool {
	rootCAPool := x509.NewCertPool()
	for _, rootCert := range certs {
		rootCAPool.AddCert(rootCert)
	}

	return rootCAPool
}

func TestNewMATLSClient(t *testing.T) {
	// Bootstrap the tests with required keys certificates
	rootCAs, err := rootCASFromFile("../../testdata/client-sample-root-ca.pem")
	if err != nil {
		t.Fatalf("load root CAs from file: %s", err.Error())
	}
	rootCAPool := rootCAPoolFromCerts(rootCAs)

	clientCerts, err := clientCertsFromFile(
		"../../testdata/client-sample-key.key", "../../testdata/client-sample-cert.pem")
	if err != nil {
		t.Fatalf("Create client certs: %s", err.Error())
	}

	config := MATLSConfig{
		InsecureSkipVerify: true,
		ClientCerts:        clientCerts,
		RootCAs:            rootCAs,
		TLSMinVersion:      tls.VersionTLS12,
	}
	wantClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
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
