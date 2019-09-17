package main

import (
	"bufio"
	"flag"
	"fmt"
	http2 "net/http"
	"os"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/cmd/config"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/http"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/version"
)

func main() {
	fmt.Println("Dynamic Client Registration Conformance Tool cli")

	flags := mustParseFlags()

	if flags.versionCmd {
		versionCmd()
	}

	runCmd(flags)
}

func versionCmd() {
	err := version.Print(bufio.NewWriter(os.Stdout))
	exitOnError(err)
	os.Exit(0)
}

func runCmd(flags flags) {
	if flags.configFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(flags.configFilePath)
	exitOnError(err)

	client := &http2.Client{Timeout: time.Second * 2}
	openIDConfig, err := openid.Get(cfg.WellknownEndpoint, client)
	exitOnError(err)

	authoriserBuilder := auth.NewAuthoriserBuilder().
		WithOpenIDConfig(openIDConfig).
		WithSSA(cfg.SSA).
		WithKID(cfg.Kid).
		WithClientID(cfg.ClientId).
		WithRedirectURIs(cfg.RedirectURIs).
		WithPrivateKey(cfg.PrivateKeyBytes).
		WithJwtExpiration(time.Hour)
	securedClient, err := http.NewBuilder().
		WithRootCAs(cfg.TransportRootCAs).
		WithTransportKeyPair(cfg.TransportCert, cfg.TransportKey).
		Build()
	exitOnError(err)

	dcr32Cfg := compliant.NewDCR32Config(
		openIDConfig,
		cfg.SSA,
		cfg.Kid,
		cfg.ClientId,
		cfg.RedirectURIs,
		cfg.PrivateKeyBytes,
	)
	scenarios := compliant.NewDCR32(dcr32Cfg, securedClient, authoriserBuilder)
	tester := compliant.NewTester(flags.filterExpression, flags.debug)

	passes, err := tester.Compliant(scenarios)
	exitOnError(err)

	if !passes {
		os.Exit(1)
	}
}

type flags struct {
	versionCmd       bool
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

	return flags{
		configFilePath:   configFilePath,
		filterExpression: filterExpression,
		debug:            debug,
		versionCmd:       versionFlag,
	}
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
