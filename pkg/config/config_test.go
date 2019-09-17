package config

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseConfig_Succeeds_WithValidConfig(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	keyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	pkeyJson, err := json.Marshal(string(keyPem))
	assert.NoError(t, err)
	json := fmt.Sprintf(`{
		"wellknown_endpoint": "https://ob19-auth1-ui.o3bank.co.uk/.well-known/openid-configuration",
   		"ssa": "ssa",
    	"kid": "kid",
    	"redirect_uris": ["https://0.0.0.0:8443/conformancesuite/callback"],
    	"client_id": "clientid",
		"private_key": %s,
    	"transport_root_cas": [
      		"-----BEGIN CERTIFICATE-----\ntransportroot1-----END CERTIFICATE-----\n",
      		"-----BEGIN CERTIFICATE-----\ntransportroot2-----END CERTIFICATE-----\n"
    	],
    	"transport_cert": "-----BEGIN CERTIFICATE-----\n\n-----END CERTIFICATE-----",
		"transport_key": "-----BEGIN CERTIFICATE-----\n\n-----END CERTIFICATE-----"}`, pkeyJson)
	cfg, err := parseConfig(bytes.NewReader([]byte(json)))
	assert.NoError(t, err)
	assert.Equal(t, "https://ob19-auth1-ui.o3bank.co.uk/.well-known/openid-configuration", cfg.WellknownEndpoint)
	assert.Equal(t, []string{"https://0.0.0.0:8443/conformancesuite/callback"}, cfg.RedirectURIs)
	assert.Equal(t, "ssa", cfg.SSA)
	assert.Equal(t, "kid", cfg.Kid)
	assert.Equal(t, "clientid", cfg.ClientId)
	assert.Equal(t, "-----BEGIN CERTIFICATE-----\n\n-----END CERTIFICATE-----", cfg.TransportCert)
	assert.Equal(t, "-----BEGIN CERTIFICATE-----\n\n-----END CERTIFICATE-----", cfg.TransportKey)
}

func Test_ParseConfig_Fails_WithInvalidPrivateKey(t *testing.T) {
	_, err := parseConfig(bytes.NewReader([]byte(`{
		"wellknown_endpoint": "https://ob19-auth1-ui.o3bank.co.uk/.well-known/openid-configuration",
   		"ssa": "ssa",
    	"kid": "kid",
		"redirect_uris": ["https://0.0.0.0:8443/conformancesuite/callback"],
		"private_key": "foobar"
	}`)))
	assert.Equal(t, "unable to parse private key bytes: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key", err.Error())
}

func Test_ParseConfig_Fails_WithInvalidConfig(t *testing.T) {
	_, err := parseConfig(bytes.NewReader([]byte(`foobar`)))
	assert.Equal(t, "unable to json decode file contents: invalid character 'o' in literal false (expecting 'a')", err.Error())
}
