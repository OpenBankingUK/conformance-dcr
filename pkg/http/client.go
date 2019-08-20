package http

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type MATLSConfig struct {
	ClientCerts        []tls.Certificate
	InsecureSkipVerify bool
	RootCAs            []*x509.Certificate
	TLSMinVersion      uint16
}

// NewMATLSClient creates a new http client that is configured for Mutually Authenticated TLS. `insecureSkipVerify`
// can be set to true if host certificates are not to be validated against a local trusted list of CA certificates.
// `caCerts` is an optional list of root CA certificates that can be used to validate host certificates.
// Can be set to nil if not required.
func NewMATLSClient(config MATLSConfig) (*http.Client, error) {
	if config.InsecureSkipVerify {
		return nil, errors.New("insecure skip verify not implemented")
	}

	tlsConfig := &tls.Config{
		Certificates:  config.ClientCerts,
		MinVersion:    config.TLSMinVersion,
		Renegotiation: tls.RenegotiateFreelyAsClient,
	}

	if config.RootCAs != nil {
		caCrtPool := x509.NewCertPool()
		for _, cert := range config.RootCAs {
			caCrtPool.AddCert(cert)
		}
		tlsConfig.RootCAs = caCrtPool
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return &http.Client{Transport: transport}, nil
}

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
