package http

import (
	"crypto/tls"
	"crypto/x509"
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
	tlsConfig := &tls.Config{
		Certificates:       config.ClientCerts,
		InsecureSkipVerify: config.InsecureSkipVerify,
		MinVersion:         config.TLSMinVersion,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
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
