package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerboseTester_Compliant(t *testing.T) {
	scenarios := Scenarios{
		NewBuilder("Scenario with one test", "Spec Link").
			TestCase(
				NewTestCaseBuilder("Always fail test").
					Step(failStep{}).
					Build(),
			).
			Build(),
	}
	tester := NewColourTester(false)

	isCompliant, err := tester.Compliant(scenarios)

	assert.NoError(t, err)
	assert.False(t, isCompliant)
}

type failStep struct{}

func (s failStep) Run(ctx step.Context) step.Result {
	return step.NewFailResult("test name", "reasons")
}
