package main

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/schema"
	"bufio"
	"flag"
	"fmt"
	http2 "net/http"
	"os"
	"strings"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/cmd/config"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/http"

	ver "bitbucket.org/openbankingteam/conformance-dcr/pkg/version"
)

func main() {
	fmt.Println("Dynamic Client Registration Conformance Tool cli")

	flags := mustParseFlags()

	vInfo := VersionInfo{
		version:    version,
		buildTime:  buildTime,
		commitHash: commitHash,
	}

	if flags.versionCmd {
		versionCmd(vInfo)
	}

	updateCheckCmd(vInfo)

	runCmd(flags)
}

func versionCmd(v VersionInfo) {
	err := v.Print(bufio.NewWriter(os.Stdout))
	exitOnError(err)
	os.Exit(0)
}

func updateCheckCmd(v VersionInfo) {
	// Check for updates and print message
	bitbucketTagsEndpoint := "https://api.bitbucket.org/2.0/repositories/openbankingteam/conformance-dcr/refs/tags"
	updMessage := getUpdateMessage(v, bitbucketTagsEndpoint)
	if updMessage != "" {
		fmt.Println(updMessage)
	}
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
		WithSoftwareID(cfg.SoftwareID).
		WithRedirectURIs(cfg.RedirectURIs).
		WithPrivateKey(cfg.PrivateKey).
		WithJwtExpiration(time.Hour)
	securedClient, err := http.NewBuilder().
		WithRootCAs(cfg.TransportRootCAs).
		WithTransportKeyPair(cfg.TransportCert, cfg.TransportKey).
		Build()
	exitOnError(err)

	const responseSchemaVersion = "3.2"
	validator, err := schema.NewValidator(responseSchemaVersion)
	exitOnError(err)

	dcr32Cfg := compliant.NewDCR32Config(
		openIDConfig,
		cfg.SSA,
		cfg.Kid,
		cfg.RedirectURIs,
		cfg.PrivateKey,
	)
	scenarios := compliant.NewDCR32(dcr32Cfg, securedClient, authoriserBuilder, validator)
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

// getUpdateMessage checks if there is an update available to the current software. An appropriate message is returned
// in both cases of either update being available or not.
func getUpdateMessage(v VersionInfo, bitbucketTagsEndpoint string) string {
	vc := ver.NewBitBucket(bitbucketTagsEndpoint)
	update, err := vc.UpdateAvailable(v.version)
	if err != nil {
		return fmt.Sprintf("Error checking for updates: %s", err.Error())
	}
	if update {
		sb := strings.Builder{}
		updMsg := fmt.Sprintf("Version %s of the this tool is out of date. Please consider updating.\n", v.version)
		sb.WriteString(updMsg)
		sb.WriteString("Please see the following URL more information:\n")
		sb.WriteString("https://bitbucket.org/openbankingteam/conformance-dcr/src/develop/README.md")
		return sb.String()
	}

	return ""
}
