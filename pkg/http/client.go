package http

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/pkg/errors"
	"net/http"
)

type MATLSClient struct {
	http.Client
}

// NewMATLSClient creates a new http client that is configured for Mutually Authenticated TLS. `insecureSkipVerify` can be set
// to true if host certificates are not to be validated against a local trusted list of CA certificates. `caCerts` is an
// optional list of root CA certificates that can be used to validate host certificates. Can be set to nil if not required.
func NewMATLSClient(certPEMBlock, keyPEMBlock []byte, insecureSkipVerify bool, caCerts []byte) (MATLSClient, error) {
	crt, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return MATLSClient{}, errors.Wrap(err, "tls.X509KeyPair")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{crt},
		InsecureSkipVerify: insecureSkipVerify,
		MinVersion:         tls.VersionTLS12,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
	}

	if caCerts != nil {
		caCrtPool := x509.NewCertPool()
		caCrtPool.AppendCertsFromPEM(caCerts)
		tlsConfig.RootCAs = caCrtPool
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	result := MATLSClient{}
	result.Client = http.Client{Transport: transport}

	return result, nil
}