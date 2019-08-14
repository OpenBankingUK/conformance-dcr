package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/config"
)

func loadConfig() config.Config {
	var configFilePath string
	flag.StringVar(&configFilePath, "config-path", "", "Config file path")
	flag.Parse()
	if configFilePath == "" {
		log.Fatalf("config-path cannot be empty")
	}
	f, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("unable to open config file %s, %v", configFilePath, err)
	}
	defer f.Close()
	rawCfg, err := ioutil.ReadAll(f)
	var cfg config.Config
	if err := json.NewDecoder(bytes.NewBuffer(rawCfg)).Decode(&cfg); err != nil {
		log.Fatalf("unable to json decode file contents, %v", err)
	}
	return cfg
}

func main() {
	fmt.Println("Dynamic Client Registration Conformance Tool cli")
	cfg := loadConfig()
	scenarios := compliant.NewDCR31(cfg.WellknownEndpoint)
	tester := compliant.NewVerboseTester()

	passes := tester.Compliant(scenarios)

	if !passes {
		fmt.Println("FAIL")
		os.Exit(1)
	}
	fmt.Println("PASS")
}
