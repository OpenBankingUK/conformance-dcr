package http

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
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
	// nolint:gosec
	tlsConfig := &tls.Config{
		Certificates:       config.ClientCerts,
		MinVersion:         config.TLSMinVersion,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
		InsecureSkipVerify: config.InsecureSkipVerify,
	}

	if config.RootCAs != nil {
		tlsConfig.RootCAs = RootCAPoolFromCerts(config.RootCAs)
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return &http.Client{Transport: transport, Timeout: time.Second * 10}, nil
}

func TlsClientCert(certPEMBlock, keyPEMBlock []byte) ([]tls.Certificate, error) {
	crt, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, errors.Wrap(err, "parse x509 key pair")
	}
	return []tls.Certificate{crt}, nil
}

func TlsCertFromFile(keyPath, certPath string) ([]tls.Certificate, error) {
	keyBlock, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, errors.Wrap(err, "tls cert from file")
	}
	certBlock, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, errors.Wrap(err, "tls cert from file")
	}

	return TlsClientCert(certBlock, keyBlock)
}

func RootCASCertificate(pemBytes []byte) (*x509.Certificate, error) {
	// Currently only support one block
	var block *pem.Block
	block, _ = pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("could not find a PEM formatted block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parse x509 certificate")
	}

	return cert, nil
}

func RootCASFromFile(path string) (*x509.Certificate, error) {
	pemBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "rootCAs from file")
	}

	return RootCASCertificate(pemBytes)
}

func RootCAs(cas []string) ([]*x509.Certificate, error) {
	rootCAs := make([]*x509.Certificate, len(cas))
	for key, rootCA := range cas {
		ca, err := RootCASCertificate([]byte(rootCA))
		if err != nil {
			return nil, errors.Wrapf(err, "building rootCAs certificate: %d", key)
		}
		rootCAs[key] = ca
	}
	return rootCAs, nil
}

func RootCAPoolFromCerts(certs []*x509.Certificate) *x509.CertPool {
	rootCAPool := x509.NewCertPool()
	for _, rootCert := range certs {
		rootCAPool.AddCert(rootCert)
	}
	return rootCAPool
}
