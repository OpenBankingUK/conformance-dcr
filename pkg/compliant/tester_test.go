package compliant

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
)

func testPassScenarios() Scenarios {
	return Scenarios{
		NewBuilder("Scenario with one test").
			TestCase(
				NewTestCaseBuilder("Always pass test").
					Step(step.NewAlwaysPass()).
					Build(),
			).
			Build(),
	}
}

func ExampleTester_Compliant() {
	tester := NewTester()

	isCompliant := tester.Compliant(testPassScenarios())

	compliantText := map[bool]string{
		false: "NOT compliant",
		true:  "compliant",
	}

	fmt.Println("Scenario with one test is " + compliantText[isCompliant])
	// Output:
	// Scenario with one test is compliant
}

func TestVerboseTester_Compliant(t *testing.T) {
	scenarios := Scenarios{
		NewBuilder("Scenario with one test").
			TestCase(
				NewTestCaseBuilder("Always fail test").
					Step(step.NewAlwaysFail()).
					Build(),
			).
			Build(),
	}
	tester := NewVerboseColourTester(false)

	isCompliant := tester.Compliant(scenarios)

	assert.False(t, isCompliant)
}
