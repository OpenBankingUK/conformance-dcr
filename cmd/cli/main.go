package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	http2 "net/http"
	"os"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/http"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/version"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant"
)

func main() {
	fmt.Println("Dynamic Client Registration Conformance Tool cli")

	flags := mustParseFlags()

	cfg, err := loadConfig(flags.configFilePath)
	exitOnError(err)

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKey))
	exitOnError(err)

	client := &http2.Client{Timeout: time.Second * 2}
	openIDConfig, err := openid.Get(cfg.WellknownEndpoint, client)
	exitOnError(err)

	authoriser := auth.NewAuthoriser(openIDConfig, cfg.SSA, cfg.Kid, cfg.ClientId, cfg.RedirectURIs, privateKey)

	securedClient, err := http.NewBuilder().
		WithRootCAs(cfg.TransportRootCAs).
		WithTransportKeyPair(cfg.TransportCert, cfg.TransportKey).
		Build()
	exitOnError(err)

	scenarios := compliant.NewDCR32(openIDConfig, securedClient, authoriser)
	tester := compliant.NewTester(flags.filterExpression, flags.debug)

	passes, err := tester.Compliant(scenarios)
	exitOnError(err)

	if !passes {
		os.Exit(1)
	}
}

type flags struct {
	configFilePath   string
	filterExpression string
	debug            bool
}

func mustParseFlags() flags {
	var configFilePath, filterExpression string
	var debug, versionFlag bool
	flag.StringVar(&configFilePath, "config-path", "", "Config file path")
	flag.StringVar(&filterExpression, "filter", "", "Filter scenarios containing value")
	flag.BoolVar(&debug, "debug", false, "Enable debug defaults to disabled")
	flag.BoolVar(&versionFlag, "version", false, "Print the version details of conformance-dcr")

	flag.Parse()

	if versionFlag {
		err := version.Print(bufio.NewWriter(os.Stdout))
		exitOnError(err)
		os.Exit(0)
	}

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
	Kid               string   `json:"kid"`
	RedirectURIs      []string `json:"redirect_uris"`
	ClientId          string   `json:"client_id"`
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

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
