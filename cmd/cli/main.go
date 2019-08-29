package main

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/http"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"io/ioutil"
	http2 "net/http"
	"os"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant"
)

func main() {
	fmt.Println("Dynamic Client Registration Conformance Tool cli")

	flags := mustParseFlags()

	cfg, err := loadConfig(flags.configFilePath)
	if err != nil {
		exitErr(err.Error())
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKey))
	if err != nil {
		exitErr(err.Error())
	}

	client := &http2.Client{Timeout: time.Second * 2}
	openIdConfig, err := openid.Get(cfg.WellknownEndpoint, client)
	if err != nil {
		exitErr(err.Error())
	}

	authoriser := auth.NewAuthoriser(openIdConfig, privateKey, cfg.SSA)

	securedClient, err := http.NewBuilder().
		WithRootCAs(cfg.TransportRootCAs).
		WithTransportKeyPair(cfg.TransportCert, cfg.TransportKey).
		Build()
	if err != nil {
		exitErr(err.Error())
	}

	scenarios := compliant.NewDCR32(cfg.WellknownEndpoint, openIdConfig.RegistrationEndpoint, securedClient, authoriser)
	tester := compliant.NewTester(flags.filterExpression, flags.debug)

	passes := tester.Compliant(scenarios)
	if !passes {
		exitErr("FAIL")
	}
	fmt.Println("PASS")
}

type flags struct {
	configFilePath   string
	filterExpression string
	debug            bool
}

func mustParseFlags() flags {
	var configFilePath, filterExpression string
	var debug bool
	flag.StringVar(&configFilePath, "config-path", "", "Config file path")
	flag.StringVar(&filterExpression, "filter", "", "Filter scenarios containing value")
	flag.BoolVar(&debug, "debug", false, "Enable debug defaults to disabled")
	flag.Parse()
	if configFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}
	return flags{
		configFilePath:   configFilePath,
		filterExpression: filterExpression,
		debug:            debug,
	}
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

func exitErr(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
