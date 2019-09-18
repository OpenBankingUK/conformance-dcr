package config

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type Config struct {
	WellknownEndpoint string          `json:"wellknown_endpoint"`
	SSA               string          `json:"ssa"`
	Kid               string          `json:"kid"`
	RedirectURIs      []string        `json:"redirect_uris"`
	ClientId          string          `json:"client_id"`
	PrivateKeyPEM     string          `json:"private_key"`
	PrivateKey        *rsa.PrivateKey `json:"-"`
	TransportRootCAs  []string        `json:"transport_root_cas"`
	TransportCert     string          `json:"transport_cert"`
	TransportKey      string          `json:"transport_key"`
}

func LoadConfig(configFilePath string) (Config, error) {
	f, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, errors.Wrapf(err, "unable to open config file %s", configFilePath)
	}
	defer f.Close()
	return parseConfig(f)
}

func parseConfig(f io.Reader) (Config, error) {
	var cfg Config
	rawCfg, err := ioutil.ReadAll(f)
	if err != nil {
		return cfg, errors.Wrap(err, "unable to read config file contents")
	}
	if err = json.NewDecoder(bytes.NewBuffer(rawCfg)).Decode(&cfg); err != nil {
		return cfg, errors.Wrap(err, "unable to json decode file contents")
	}
	privateKeyBytes, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKeyPEM))
	if err != nil {
		return cfg, errors.Wrap(err, "unable to parse private key bytes")
	}
	cfg.PrivateKey = privateKeyBytes
	return cfg, nil
}
