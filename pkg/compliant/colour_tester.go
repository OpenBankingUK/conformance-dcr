package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"fmt"
	"github.com/logrusorgru/aurora"
)

func NewColourTester(debug bool) Tester {
	return colourTester{
		debug: debug,
	}
}

type colourTester struct {
	debug bool
}

func (t colourTester) Compliant(scenarios Scenarios) (bool, error) {
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
	return ok, nil
}

func (t colourTester) printColourTestResult(result step.Result) {
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
