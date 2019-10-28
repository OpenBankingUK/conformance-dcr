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
    	"software_id": "softwareId",
		"private_key": %s,
    	"transport_root_cas": [
      		"-----BEGIN CERTIFICATE-----\ntransportroot1-----END CERTIFICATE-----\n",
      		"-----BEGIN CERTIFICATE-----\ntransportroot2-----END CERTIFICATE-----\n"
    	],
    	"transport_cert": "-----BEGIN CERTIFICATE-----\n\n-----END CERTIFICATE-----",
		"transport_key": "-----BEGIN CERTIFICATE-----\n\n-----END CERTIFICATE-----",
		"get_implemented": true,
		"put_implemented": false,
		"delete_implemented": true
	}`, pkeyJson)
	cfg, err := parseConfig(bytes.NewReader([]byte(json)))
	expectedCfg := Config{
		WellknownEndpoint: "https://ob19-auth1-ui.o3bank.co.uk/.well-known/openid-configuration",
		SSA:               "ssa",
		Kid:               "kid",
		RedirectURIs:      []string{"https://0.0.0.0:8443/conformancesuite/callback"},
		SoftwareID:        "softwareId",
		PrivateKeyPEM:     string(keyPem),
		PrivateKey:        key,
		TransportRootCAs: []string{
			"-----BEGIN CERTIFICATE-----\ntransportroot1-----END CERTIFICATE-----\n",
			"-----BEGIN CERTIFICATE-----\ntransportroot2-----END CERTIFICATE-----\n",
		},
		TransportCert:     "-----BEGIN CERTIFICATE-----\n\n-----END CERTIFICATE-----",
		TransportKey:      "-----BEGIN CERTIFICATE-----\n\n-----END CERTIFICATE-----",
		GetImplemented:    true,
		PutImplemented:    false,
		DeleteImplemented: true,
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedCfg, cfg)
}

func Test_ParseConfig_Fails_WithInvalidPrivateKey(t *testing.T) {
	_, err := parseConfig(bytes.NewReader([]byte(`{
		"wellknown_endpoint": "https://ob19-auth1-ui.o3bank.co.uk/.well-known/openid-configuration",
   		"ssa": "ssa",
    	"kid": "kid",
		"redirect_uris": ["https://0.0.0.0:8443/conformancesuite/callback"],
		"private_key": "foobar"
	}`)))
	assert.EqualError(
		t,
		err,
		"unable to parse private key bytes: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key",
	)
}

func Test_ParseConfig_Fails_WithInvalidConfig(t *testing.T) {
	_, err := parseConfig(bytes.NewReader([]byte(`foobar`)))
	assert.EqualError(
		t,
		err,
		"unable to json decode file contents: invalid character 'o' in literal false (expecting 'a')",
	)
}
