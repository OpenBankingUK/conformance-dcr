package main

import (
	"bufio"
	"flag"
	"fmt"
	http2 "net/http"
	"os"
	"strings"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
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

	cfg, err := LoadConfig(flags.configFilePath)
	exitOnError(err)

	client := &http2.Client{Timeout: time.Second * 2}
	openIDConfig, err := openid.Get(cfg.WellknownEndpoint, client)
	exitOnError(err)

	dcr32Cfg, err := compliant.NewDCR32Config(
		openIDConfig,
		cfg.SSA,
		cfg.Kid,
		cfg.SoftwareStatementId,
		cfg.RedirectURIs,
		cfg.SigningKeyPEM,
		cfg.TransportKeyPEM,
		cfg.TransportCertPEM,
		cfg.TransportRootCAsPEM,
		cfg.GetImplemented,
		cfg.PutImplemented,
		cfg.DeleteImplemented,
		flags.tokenEndpointRS256Method,
	)
	exitOnError(err)

	manifest, err := compliant.NewDCR32(dcr32Cfg)
	exitOnError(err)

	if flags.filterExpression != "" {
		manifest, err = compliant.NewFilteredManifest(manifest, flags.filterExpression)
		exitOnError(err)
	}

	tester := compliant.NewTester()

	printer := compliant.NewPrinter(flags.debug)
	tester.AddListener(printer.Print)

	if flags.report {
		reporterFunc := compliant.NewReporter(flags.debug, "report.json")
		tester.AddListener(reporterFunc.Report)
	}

	passes, err := tester.Compliant(manifest)
	exitOnError(err)

	if !passes {
		os.Exit(1)
	}
}

type flags struct {
	versionCmd               bool
	configFilePath           string
	filterExpression         string
	debug                    bool
	report                   bool
	tokenEndpointRS256Method bool
}

func mustParseFlags() flags {
	var configFilePath, filterExpression string
	var debug, report, versionFlag, tokenEndpointRS256Method bool
	flag.StringVar(&configFilePath, "config-path", "", "Config file path")
	flag.StringVar(&filterExpression, "filter", "", "Filter scenarios containing value")
	flag.BoolVar(&debug, "debug", false, "Enable debug defaults to disabled")
	flag.BoolVar(&report, "report", false, "Enable report output defaults to disabled")
	flag.BoolVar(&versionFlag, "version", false, "Print the version details of conformance-dcr")
	flag.BoolVar(&tokenEndpointRS256Method, "rs256", false, "Run test suite with RS256 (testing only)")
	flag.Parse()

	return flags{
		configFilePath:           configFilePath,
		filterExpression:         filterExpression,
		debug:                    debug,
		report:                   report,
		versionCmd:               versionFlag,
		tokenEndpointRS256Method: tokenEndpointRS256Method,
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
		return fmt.Sprintf("error checking for updates: %s", err.Error())
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
