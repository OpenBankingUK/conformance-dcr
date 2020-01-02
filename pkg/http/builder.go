package http

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"net/http"
)

type MATLSClientBuilder interface {
	WithRootCAs(rootCAs []string) MATLSClientBuilder
	WithTransportKeyPair(certPEMBlock, keyPEMBlock string) MATLSClientBuilder
	Build() (*http.Client, error)
}

type mTLSClientBuilder struct {
	certPEMBlock, keyPEMBlock *string
	rootCAs                   *[]string
	tlsSkipVerify             bool
}

func NewBuilder() *mTLSClientBuilder {
	return &mTLSClientBuilder{
		certPEMBlock:  nil,
		keyPEMBlock:   nil,
		rootCAs:       nil,
		tlsSkipVerify: false,
	}
}

func (b *mTLSClientBuilder) WithRootCAs(rootCAs []string) *mTLSClientBuilder {
	b.rootCAs = &rootCAs
	return b
}

func (b *mTLSClientBuilder) WithTlsSkipVerify(tlsSkipVerify bool) *mTLSClientBuilder {
	b.tlsSkipVerify = tlsSkipVerify
	return b
}

func (b *mTLSClientBuilder) WithTransportKeyPair(certPEMBlock, keyPEMBlock string) *mTLSClientBuilder {
	b.certPEMBlock = &certPEMBlock
	b.keyPEMBlock = &keyPEMBlock
	return b
}

func (b *mTLSClientBuilder) Build() (*http.Client, error) {
	if b.certPEMBlock == nil || b.keyPEMBlock == nil {
		return nil, errors.New("can't build a mtls client without cert and key")
	}

	clientCerts, err := TlsClientCert([]byte(*b.certPEMBlock), []byte(*b.keyPEMBlock))
	if err != nil {
		return nil, errors.Wrap(err, "building mTLS http client")
	}

	if b.rootCAs == nil {
		return nil, errors.New("can't build a mTLS client without rootCAs")
	}

	rootCAs, err := RootCAs(*b.rootCAs)
	if err != nil {
		return nil, errors.Wrap(err, "building mTLS http client")
	}

	config := MATLSConfig{
		ClientCerts:        clientCerts,
		InsecureSkipVerify: b.tlsSkipVerify,
		RootCAs:            rootCAs,
		TLSMinVersion:      tls.VersionTLS12,
	}

	return NewMATLSClient(config)
}
