package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/server"
	"github.com/labstack/echo"
	"github.com/spf13/cobra"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/resty.v1"
)

const (
	certFile = "./configs/dcr.crt"
	keyFile  = "./configs/dcr.key"
	version  = "v1.0.0"
	host     = "127.0.0.1"
)

type ConfigBuilder struct {
	logger *logrus.Logger
}

func NewConfigBuilder(logger *logrus.Logger) ConfigBuilder {
	return ConfigBuilder{logger}
}

func (cfgBuilder ConfigBuilder) InitConfig() {
	cfgBuilder.logger.SetNoLock()
	cfgBuilder.logger.SetFormatter(&prefixed.TextFormatter{
		DisableColors:    false,
		ForceColors:      true,
		TimestampFormat:  time.RFC3339,
		FullTimestamp:    true,
		DisableTimestamp: false,
		ForceFormatting:  true,
	})
	level, err := logrus.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		cfgBuilder.PrintConfigFlags()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cfgBuilder.logger.SetLevel(level)
	if viper.GetBool("log_http_file") {
		httpLogFile, err := os.OpenFile("http-trace.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			logrus.Warn("cannot set http trace file")
		} else {
			resty.SetLogger(httpLogFile)
		}
	}

	resty.SetDebug(viper.GetBool("log_http_trace"))

	cfgBuilder.PrintConfigFlags()
}

func (cfgBuilder ConfigBuilder) PrintConfigFlags() {
	cfgBuilder.logger.WithFields(logrus.Fields{
		"log_level":      viper.GetString("log_level"),
		"log_tracer":     viper.GetBool("log_tracer"),
		"log_http_trace": viper.GetBool("log_http_trace"),
		"log_http_file":  viper.GetBool("log_http_file"),
		"log_to_file":    viper.GetBool("log_to_file"),
		"port":           viper.GetInt("port"),
	}).Info("configuration flags")
}

func initRootCmd() *cobra.Command {
	configBuilder := NewConfigBuilder(logrus.StandardLogger())
	rootCmd := &cobra.Command{
		Use:   "dcr_server",
		Short: "Dynamic Client Registration Conformance Suite",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := configBuilder.logger.WithField("app", "server")
			server := server.NewServer(echo.New(), logger, version)
			address := fmt.Sprintf("%s:%d", host, viper.GetInt("port"))
			logger.Infof("listening on https://%s", address)
			return server.StartTLS(address, certFile, keyFile)
		},
	}
	rootCmd.PersistentFlags().String("log_level", "INFO", "Log level")
	rootCmd.PersistentFlags().Bool("log_tracer", false, "Enable tracer logging")
	rootCmd.PersistentFlags().Bool("log_http_trace", false, "Enable HTTP logging")
	rootCmd.PersistentFlags().Int("port", 8443, "Server port")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	cobra.OnInitialize(configBuilder.InitConfig)

	return rootCmd
}

func main() {
	rootCmd := initRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
