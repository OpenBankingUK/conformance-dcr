package main

import (
	"bufio"
	"crypto/rsa"
	"flag"
	"fmt"
	http2 "net/http"
	"os"
	"strings"
	"time"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/openid"
	ver "github.com/OpenBankingUK/conformance-dcr/pkg/version"
	"github.com/dgrijalva/jwt-go"
)

func main() {
	fmt.Println("Dynamic Client Registration Conformance Tool cli")

	patchJwtLibraryBug()

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

	client := &http2.Client{Timeout: time.Second * 5}
	openIDConfig, err := openid.Get(cfg.WellknownEndpoint, client)
	exitOnError(err)

	dcr32Cfg, err := compliant.NewDCR32Config(
		openIDConfig,
		cfg.SSA,
		cfg.Aud,
		cfg.Kid,
		cfg.Issuer,
		cfg.RedirectURIs,
		cfg.SigningKeyPEM,
		cfg.TransportKeyPEM,
		cfg.TransportCertPEM,
		cfg.TransportCertSubjectDN,
		cfg.TransportRootCAsPEM,
		cfg.GetImplemented,
		cfg.PutImplemented,
		cfg.DeleteImplemented,
		flags.tlsSkipVerify,
		cfg.SpecVersion,
		cfg.SSAs,
	)
	exitOnError(err)

	manifest, err := compliant.NewSpecManifest(cfg.SpecVersion, dcr32Cfg)
	exitOnError(err)

	if flags.filterExpression != "" {
		manifest, err = compliant.NewFilteredManifest(manifest, flags.filterExpression)
		exitOnError(err)
	}

	tester := compliant.NewTester()

	printer := compliant.NewPrinter(flags.debug)
	tester.AddListener(printer.Print)

	doneSignal := make(chan bool)
	serverAddr := serverAddress(flags.httpServerPort)
	if flags.report {
		reporterFunc := compliant.NewReporter(runConfig(cfg), flags.debug, doneSignal, serverAddr)
		tester.AddListener(reporterFunc.Report)
	}

	passes, err := tester.Compliant(manifest)
	exitOnError(err)

	if flags.report {
		waitForDownloadOrTimeout(serverAddr, doneSignal)
	}

	if !passes {
		os.Exit(1)
	}
}

func runConfig(config Config) compliant.RunConfig {
	return compliant.RunConfig{
		WellknownEndpoint: config.WellknownEndpoint,
		GetImplemented:    config.GetImplemented,
		PutImplemented:    config.PutImplemented,
		DeleteImplemented: config.DeleteImplemented,
		Environment:       config.Environment,
		Brand:             config.Brand,
	}
}

func waitForDownloadOrTimeout(serverAddr string, doneSignal <-chan bool) {
	fmt.Printf("To download report open webpage http://%s\n", serverAddr)
	fmt.Println("Waiting for report download...")
	select {
	case <-doneSignal:
		fmt.Println("Report download completed, exiting.")
		break
	case <-time.After(2 * time.Minute):
		fmt.Println("Time out waiting for report download, exiting.")
		break
	}
}

func serverAddress(port string) string {
	return fmt.Sprintf("0.0.0.0:%s", port)
}

type flags struct {
	versionCmd       bool
	configFilePath   string
	filterExpression string
	debug            bool
	report           bool
	tlsSkipVerify    bool
	httpServerPort   string
}

func mustParseFlags() flags {
	var configFilePath, filterExpression, httpServerPort string
	var debug, report, versionFlag, tlsSkipVerify bool
	flag.StringVar(&configFilePath, "config-path", "", "Config file path")
	flag.StringVar(&filterExpression, "filter", "", "Filter scenarios containing value")
	flag.StringVar(&httpServerPort, "port", "8080", "Http server port for report download")
	flag.BoolVar(&debug, "debug", false, "Enable debug defaults to disabled")
	flag.BoolVar(&report, "report", false, "Enable report output defaults to disabled")
	flag.BoolVar(&versionFlag, "version", false, "Print the version details of conformance-dcr")
	flag.BoolVar(&tlsSkipVerify, "tlsskipverify", false, "Skip ssl cert verify")
	flag.Parse()

	return flags{
		configFilePath:   configFilePath,
		filterExpression: filterExpression,
		debug:            debug,
		report:           report,
		versionCmd:       versionFlag,
		tlsSkipVerify:    tlsSkipVerify,
		httpServerPort:   httpServerPort,
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
	vc := ver.NewGitHub(bitbucketTagsEndpoint)
	update, err := vc.UpdateAvailable(v.version)
	if err != nil {
		return fmt.Sprintf("error checking for updates: %s", err.Error())
	}
	if update {
		sb := strings.Builder{}
		updMsg := fmt.Sprintf("Version %s of the this tool is out of date. Please consider updating.\n", v.version)
		sb.WriteString(updMsg)
		sb.WriteString("Please see the following URL more information:\n")
		sb.WriteString("https://github.com/OpenBankingUK/conformance-dcr/src/develop/README.md")
		return sb.String()
	}

	return ""
}

func patchJwtLibraryBug() {
	jwt.SigningMethodPS256.Options.SaltLength = rsa.PSSSaltLengthEqualsHash
	jwt.SigningMethodPS384.Options.SaltLength = rsa.PSSSaltLengthEqualsHash
	jwt.SigningMethodPS512.Options.SaltLength = rsa.PSSSaltLengthEqualsHash
}
