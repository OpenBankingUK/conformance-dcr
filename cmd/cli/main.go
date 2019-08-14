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
)

type Config struct {
	WellknownEndpoint string `json:"wellknown_endpoint"`
	SSA               string `json:"ssa"`
}

func loadConfig() Config {
	var configFilePath string
	flag.StringVar(&configFilePath, "config-path", "", "Config file path")
	flag.Parse()
	if configFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}
	f, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("unable to open config file %s, %v", configFilePath, err)
	}
	defer f.Close()
	rawCfg, err := ioutil.ReadAll(f)
	var cfg Config
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
