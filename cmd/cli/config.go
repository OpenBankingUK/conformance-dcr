package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type Config struct {
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
}

func LoadConfig(configFilePath string) (Config, error) {
	f, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, errors.Wrap(err, "load config")
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
	return cfg, nil
}
