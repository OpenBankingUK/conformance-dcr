package main

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/http"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant"
)

func main() {
	fmt.Println("Dynamic Client Registration Conformance Tool cli")

	var configFilePath string
	flag.StringVar(&configFilePath, "config-path", "", "Config file path")
	flag.Parse()
	if configFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := loadConfig(configFilePath)
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKey))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	httpClient, err := http.NewBuilder().
		WithRootCAs(cfg.TransportRootCAs).
		WithTransportKeyPair(cfg.TransportCert, cfg.TransportKey).
		Build()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	scenarios := compliant.NewDCR32(cfg.WellknownEndpoint, cfg.SSA, privateKey, httpClient)
	tester := compliant.NewColourTester()

	passes := tester.Compliant(scenarios)
	if !passes {
		fmt.Println("FAIL")
		os.Exit(1)
	}
	fmt.Println("PASS")
}

type Config struct {
	WellknownEndpoint string   `json:"wellknown_endpoint"`
	SSA               string   `json:"ssa"`
	PrivateKey        string   `json:"private_key"`
	TransportRootCAs  []string `json:"transport_root_cas"`
	TransportCert     string   `json:"transport_cert"`
	TransportKey      string   `json:"transport_key"`
}

func loadConfig(configFilePath string) (Config, error) {
	var cfg Config
	f, err := os.Open(configFilePath)
	if err != nil {
		return cfg, errors.Wrapf(err, "unable to open config file %s", configFilePath)
	}
	defer f.Close()
	rawCfg, err := ioutil.ReadAll(f)
	if err != nil {
		return cfg, errors.Wrap(err, "unable to read config file contents")
	}
	if err := json.NewDecoder(bytes.NewBuffer(rawCfg)).Decode(&cfg); err != nil {
		return cfg, errors.Wrapf(err, "unable to json decode file contents")
	}
	return cfg, nil
}
