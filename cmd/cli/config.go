package main

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type Config struct {
	SpecVersion         string   `json:"spec_version"`
	WellknownEndpoint   string   `json:"wellknown_endpoint"`
	SSA                 string   `json:"ssa"`
	Kid                 string   `json:"kid"`
	Aud                 string   `json:"aud"`
	RedirectURIs        []string `json:"redirect_uris"`
	Issuer              string   `json:"issuer"`
	SigningKeyPEM       string   `json:"private_key"`
	TransportRootCAsPEM []string `json:"transport_root_cas"`
	TransportCertPEM    string   `json:"transport_cert"`
	TransportKeyPEM     string   `json:"transport_key"`
	GetImplemented      bool     `json:"get_implemented"`
	PutImplemented      bool     `json:"put_implemented"`
	DeleteImplemented   bool     `json:"delete_implemented"`
	Environment         string   `json:"environment"`
	Brand               string   `json:"brand"`
}

func LoadConfig(configFilePath string) (Config, error) {
	f, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, errors.Wrap(err, "load config")
	}
	defer f.Close()

	config, err := parseConfig(f)
	if err != nil {
		return Config{}, errors.Wrap(err, "load config")
	}

	err = validateConfig(config)
	if err != nil {
		return Config{}, errors.Wrap(err, "load config")
	}

	return config, nil
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
	return cfg, nil
}

func validateConfig(config Config) error {
	if !compliant.IsSupportedSpecVersion(config.SpecVersion) {
		return errors.New("missing or invalid config property Specification version `spec_version`")
	}
	if config.WellknownEndpoint == "" {
		return errors.New("missing config property Well-known Endpoint `wellknown_endpoint`")
	}
	if config.Environment == "" {
		return errors.New("missing config property Environment `environment`")
	}
	if config.Brand == "" {
		return errors.New("missing config property Brand `brand`")
	}
	return nil
}
