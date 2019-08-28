package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	tester := NewColourTester(false)

	isCompliant := tester.Compliant(scenarios)

	assert.False(t, isCompliant)
}
