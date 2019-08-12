package main

import (
	"fmt"
	"os"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant"
)

func main() {
	fmt.Println("Dynamic Client Registration Conformance Tool cli")

	scenarios := compliant.NewDCR31()
	tester := compliant.NewVerboseTester()

	passes := tester.Compliant(scenarios)

	if !passes {
		fmt.Println("FAIL")
		os.Exit(1)
	}
	fmt.Println("PASS")
}
