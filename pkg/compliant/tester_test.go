package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"fmt"
)

func testScenarios() Scenarios {
	return Scenarios{
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
}

func ExampleTester_Compliant() {
	tester := NewTester()

	isCompliant := tester.Compliant(testScenarios())

	compliantText := map[bool]string{
		false: "NOT compliant",
		true:  "compliant",
	}

	fmt.Println("Scenario with one test is " + compliantText[isCompliant])
	// Output:
	// Scenario with one test is compliant
}

func ExampleNewVerboseTester() {
	tester := NewVerboseTester()

	isCompliant := tester.Compliant(testScenarios())

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
