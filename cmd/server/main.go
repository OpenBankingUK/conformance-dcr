package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/resty.v1"
)

const (
	certFile = "./configs/dcr.crt"
	keyFile  = "./configs/dcr.key"
)

var (
	logger  = logrus.StandardLogger()
	rootCmd = &cobra.Command{
		Use:   "dcr_server",
		Short: "Dynamic Client Registration Conformance Suite",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logger.WithField("app", "server")
			server := echo.New()
			address := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
			logger.Infof("listening on https://%s", address)
			return server.StartTLS(address, certFile, keyFile)
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		fmt.Fprint(os.Stderr, "\n")
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("log_level", "INFO", "Log level")
	rootCmd.PersistentFlags().Bool("log_tracer", false, "Enable tracer logging")
	rootCmd.PersistentFlags().Bool("log_http_trace", false, "Enable HTTP logging")
	rootCmd.PersistentFlags().Int("port", 8443, "Server port")
	rootCmd.PersistentFlags().String("host", "127.0.0.1", "Server host")
	rootCmd.PersistentFlags().Bool("disable_jws", false, "Disable JWS Signature")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Fprint(os.Stderr, err)
		fmt.Fprint(os.Stderr, "\n")
		os.Exit(1)
	}

	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	logger.SetNoLock()
	logger.SetFormatter(&prefixed.TextFormatter{
		DisableColors:    false,
		ForceColors:      true,
		TimestampFormat:  time.RFC3339,
		FullTimestamp:    true,
		DisableTimestamp: false,
		ForceFormatting:  true,
	})
	level, err := logrus.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		printConfigurationFlags()
		fmt.Fprint(os.Stderr, err)
		fmt.Fprint(os.Stderr, "\n")
		os.Exit(1)
	}
	logger.SetLevel(level)

	tracer.Silent = !viper.GetBool("log_tracer")
	if viper.GetBool("log_to_file") {
		f, err := os.OpenFile("suite.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			// continue as normal
		} else {
			mw := f // io.MultiWriter(os.Stdout, f)
			logrus.SetOutput(mw)
			logger.SetFormatter(&prefixed.TextFormatter{
				DisableColors:    true,
				ForceColors:      false,
				TimestampFormat:  time.RFC3339,
				FullTimestamp:    true,
				DisableTimestamp: false,
				ForceFormatting:  true,
			})

		}
	}
	if viper.GetBool("log_http_file") {
		httpLogFile, err := os.OpenFile("http-trace.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			logrus.Warn("cannot set http trace file")
		} else {
			resty.SetLogger(httpLogFile)
		}
	}

	resty.SetDebug(viper.GetBool("log_http_trace"))

	printConfigurationFlags()
}

func printConfigurationFlags() {
	logger.WithFields(logrus.Fields{
		"log_level":      viper.GetString("log_level"),
		"log_tracer":     viper.GetBool("log_tracer"),
		"log_http_trace": viper.GetBool("log_http_trace"),
		"log_http_file":  viper.GetBool("log_http_file"),
		"log_to_file":    viper.GetBool("log_to_file"),
		"port":           viper.GetInt("port"),
		"tracer.Silent":  tracer.Silent,
	}).Info("configuration flags")
}
