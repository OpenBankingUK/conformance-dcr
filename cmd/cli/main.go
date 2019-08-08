package main

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Dynamic Client Registration Conformance Tool cli")

	scenarios := compliant.NewGoogleScenarios()
	tester := compliant.NewVerboseTester()

	passes := tester.Compliant(scenarios)

	if !passes {
		fmt.Println("FAIL")
		os.Exit(1)
	}
	fmt.Println("PASS")
}
