package http

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/pkg/errors"
	"net/http"
)

type MATLSConfig struct {
	certPEMBlock       []byte
	keyPEMBlock        []byte
	insecureSkipVerify bool
	caCerts            []byte
}

// NewMATLSClient creates a new http client that is configured for Mutually Authenticated TLS. `insecureSkipVerify`
// can be set to true if host certificates are not to be validated against a local trusted list of CA certificates.
// `caCerts` is an optional list of root CA certificates that can be used to validate host certificates.
// Can be set to nil if not required.
func NewMATLSClient(config MATLSConfig) (*http.Client, error) {
	crt, err := tls.X509KeyPair(config.certPEMBlock, config.keyPEMBlock)
	if err != nil {
		return nil, errors.Wrap(err, "tls.X509KeyPair")
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{crt},
		InsecureSkipVerify: config.insecureSkipVerify,
		MinVersion:         tls.VersionTLS12,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
	}

	if config.caCerts != nil {
		caCrtPool := x509.NewCertPool()
		caCrtPool.AppendCertsFromPEM(config.caCerts)
		tlsConfig.RootCAs = caCrtPool
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return &http.Client{Transport: transport}, nil
}
