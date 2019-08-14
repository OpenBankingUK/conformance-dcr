package main

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/config"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.AddConfigPath("configs")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("unable to read config file: %v", err)
	}
}

func main() {
	fmt.Println("Dynamic Client Registration Conformance Tool cli")
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to parse config %v\n", err)
	}
	scenarios := compliant.NewDCR31(cfg.WellknownEndpoint)
	tester := compliant.NewVerboseTester()

	passes := tester.Compliant(scenarios)

	if !passes {
		fmt.Println("FAIL")
		os.Exit(1)
	}
	fmt.Println("PASS")
}
