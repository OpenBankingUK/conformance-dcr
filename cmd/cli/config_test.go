package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

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
	// The config should contain SSA or SSAs (list of SSAs). This test is about checking the parsing.
	configJson := fmt.Sprintf(`{
		"wellknown_endpoint": "https://ob19-auth1-ui.o3bank.co.uk/.well-known/openid-configuration",
   		"ssa": "ssa",
    	"kid": "kid",
    	"aud": "aud",
    	"redirect_uris": ["https://0.0.0.0:8443/conformancesuite/callback"],
    	"issuer": "softwareId",
		"private_key": %s,
    	"transport_root_cas": [
      		"-----BEGIN CERTIFICATE-----\ntransportroot1-----END CERTIFICATE-----\n",
      		"-----BEGIN CERTIFICATE-----\ntransportroot2-----END CERTIFICATE-----\n"
    	],
    	"transport_cert": "-----BEGIN CERTIFICATE-----\n\n-----END CERTIFICATE-----",
		"transport_key": "-----BEGIN CERTIFICATE-----\n\n-----END CERTIFICATE-----",
		"get_implemented": true,
		"put_implemented": false,
		"delete_implemented": true,
		"environment": "environment",
		"brand": "brand",
		"ssas": ["ssa1", "ssa2"]
	}`, pkeyJson)
	cfg, err := parseConfig(bytes.NewReader([]byte(configJson)))
	expectedCfg := Config{
		WellknownEndpoint: "https://ob19-auth1-ui.o3bank.co.uk/.well-known/openid-configuration",
		SSA:               "ssa",
		Kid:               "kid",
		Aud:               "aud",
		RedirectURIs:      []string{"https://0.0.0.0:8443/conformancesuite/callback"},
		Issuer:            "softwareId",
		SigningKeyPEM:     string(keyPem),
		TransportRootCAsPEM: []string{
			"-----BEGIN CERTIFICATE-----\ntransportroot1-----END CERTIFICATE-----\n",
			"-----BEGIN CERTIFICATE-----\ntransportroot2-----END CERTIFICATE-----\n",
		},
		TransportCertPEM:  "-----BEGIN CERTIFICATE-----\n\n-----END CERTIFICATE-----",
		TransportKeyPEM:   "-----BEGIN CERTIFICATE-----\n\n-----END CERTIFICATE-----",
		GetImplemented:    true,
		PutImplemented:    false,
		DeleteImplemented: true,
		Environment:       "environment",
		Brand:             "brand",
		SSAs:              []string{"ssa1", "ssa2"},
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedCfg, cfg)
}

func Test_ParseConfig_Fails_WithInvalidConfig(t *testing.T) {
	_, err := parseConfig(bytes.NewReader([]byte(`foobar`)))
	assert.EqualError(
		t,
		err,
		"unable to json decode file contents: invalid character 'o' in literal false (expecting 'a')",
	)
}

func Test_ReadsConfigFromFile(t *testing.T) {
	config, err := LoadConfig("testdata/config.json.sample")
	require.NoError(t, err)

	assert.Equal(t, "https://ob19-auth1-ui.o3bank.co.uk/.well-known/openid-configuration", config.WellknownEndpoint)
	assert.Equal(t, "ssa", config.SSA)
	assert.Equal(t, "kid", config.Kid)
	assert.Equal(t, "aud", config.Aud)
	assert.Equal(t, []string{"redirecturi"}, config.RedirectURIs)
	assert.Equal(t, "private_key", config.SigningKeyPEM)
	assert.Equal(t, []string{"cert 1", "cert 2"}, config.TransportRootCAsPEM)
	assert.Equal(t, "transport cert", config.TransportCertPEM)
	assert.Equal(t, "transport key", config.TransportKeyPEM)
	assert.True(t, config.GetImplemented)
	assert.True(t, config.PutImplemented)
	assert.True(t, config.DeleteImplemented)
	assert.Equal(t, "sandbox", config.Environment)
	assert.Equal(t, "Brand/product", config.Brand)
	assert.Equal(t, []string([]string(nil)), config.SSAs)
}

func Test_ReadsConfigFromSSAsFile(t *testing.T) {
	config, err := LoadConfig("testdata/config.json.ssas.sample")
	require.NoError(t, err)

	assert.Equal(t, "", config.SSA)
	assert.Equal(t, []string{"ssa1", "ssa2"}, config.SSAs)
}

func Test_ReadsConfigFromFile_HandlesError(t *testing.T) {
	_, err := LoadConfig("non_existing_file")
	require.EqualError(t, err, "load config: open non_existing_file: no such file or directory")
}
