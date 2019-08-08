package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"fmt"
)

var scenarios = Scenarios{
	NewScenario(
		"Scenario with one test",
		[]TestCase{
			NewTestCase(
				1,
				"Always pass test",
				[]step.Step{
					step.NewAlwaysPass(1),
				},
			),
		},
	),
}

func ExampleCompliant() {
	tester := NewTester()

	isCompliant := tester.Compliant(scenarios)

	compliantText := map[bool]string{
		false: "NOT compliant",
		true:  "compliant",
	}

	fmt.Println("Scenario with one test is " + compliantText[isCompliant])
	// Output:
	// Scenario with one test is compliant
}

func ExampleVerboseCompliant() {
	tester := NewVerboseTester()

	isCompliant := tester.Compliant(scenarios)

	compliantText := map[bool]string{
		false: "FAIL",
		true:  "PASS",
	}

	fmt.Println(compliantText[isCompliant])
	// Output:
	// === Scenario: Scenario with one test
	// 	Test case: Always pass test
	// 		PASS always dumb pass step
	// PASS
}
