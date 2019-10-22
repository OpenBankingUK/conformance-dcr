package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerboseTester_Compliant(t *testing.T) {
	scenarios := Scenarios{
		NewBuilder("#1", "Scenario with one test", "Spec Link").
			TestCase(
				NewTestCaseBuilder("Always fail test").
					Step(failStep{}).
					Build(),
			).
			Build(),
	}
	manifest, err := NewManifest("test", "0.0", scenarios)
	assert.NoError(t, err)
	tester := NewColourTester(false)

	isCompliant, err := tester.Compliant(manifest)

	assert.NoError(t, err)
	assert.False(t, isCompliant)
}

type failStep struct{}

func (s failStep) Run(ctx step.Context) step.Result {
	return step.NewFailResult("test name", "reasons")
}
