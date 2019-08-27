package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"fmt"
	"github.com/logrusorgru/aurora"
)

type Tester interface {
	Compliant(scenarios Scenarios) bool
}

func NewTester() Tester {
	return tester{}
}

type tester struct{}

func (t tester) Compliant(scenarios Scenarios) bool {
	ok := true
	for _, scenarios := range scenarios {
		results := scenarios.Run()
		ok = ok && !results.Fail()
	}
	return ok
}

func NewVerboseColourTester(debug bool) Tester {
	return colourVerboseTester{
		debug: debug,
	}
}

type colourVerboseTester struct {
	debug bool
}

func (t colourVerboseTester) Compliant(scenarios Scenarios) bool {
	ok := true
	for _, scenario := range scenarios {
		scenarioResult := scenario.Run()
		fmt.Printf("=== Scenario: %s\n", scenarioResult.Name)
		for _, testCasesResult := range scenarioResult.TestCaseResults {
			fmt.Printf("\tTest case: %s\n", testCasesResult.Name)
			for _, stepResult := range testCasesResult.Results {
				t.printColourTestResult(stepResult)
			}
		}
		ok = ok && !scenarioResult.Fail()
	}
	return ok
}

func (t colourVerboseTester) printColourTestResult(result step.Result) {
	if result.Pass {
		fmt.Printf("\t\t%s %s\n", aurora.Green("PASS"), result.Name)
	} else {
		fmt.Printf(
			"\t\t%s %s: %s\n",
			aurora.Red("FAIL"),
			result.Name,
			result.FailReason,
		)
	}
	if t.debug {
		printColourDebugMessages(result.Debug)
	}
}

func printColourDebugMessages(log step.DebugMessages) {
	for _, msg := range log.Item {
		fmt.Printf(
			"%s %s\n",
			msg.Time.Format("2006/01/02 15:04:05"),
			aurora.Gray(15, msg.Message),
		)
	}
}
