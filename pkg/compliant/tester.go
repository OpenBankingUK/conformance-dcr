package compliant

import "fmt"

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

func NewVerboseTester() Tester {
	return verboseTester{}
}

type verboseTester struct{}

func (t verboseTester) Compliant(scenarios Scenarios) bool {
	ok := true
	for _, scenario := range scenarios {
		scenarioResult := scenario.Run()
		fmt.Printf("=== Scenario: %s\n", scenarioResult.Name)
		for _, testCasesResult := range scenarioResult.TestCaseResults {
			fmt.Printf("\tTest case: %s\n", testCasesResult.Name)
			for _, stepResult := range testCasesResult.Results {
				if stepResult.Pass {
					fmt.Printf("\t\tPASS %s\n", stepResult.Name)
				} else {
					fmt.Printf("\t\tFAIL %s: %s\n", stepResult.Name, stepResult.Message)
				}
			}
		}
		ok = ok && !scenarioResult.Fail()
	}
	return ok
}
