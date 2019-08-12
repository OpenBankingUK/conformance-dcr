package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTestCaseResults_Fail_FalseIfPasses(t *testing.T) {
	tcs := TestCaseResults{
		TestCaseResult{
			Results: step.Results{
				step.Result{Pass: true},
			},
		},
	}

	assert.False(t, tcs.Fail())
}

func TestTestCaseResults_Fail_TrueOneFails(t *testing.T) {
	tcs := TestCaseResults{
		TestCaseResult{
			Results: step.Results{
				step.Result{Pass: true},
			},
		},
		TestCaseResult{
			Results: step.Results{
				step.Result{Pass: false},
			},
		},
	}

	assert.True(t, tcs.Fail())
}

func TestTestCase_Run_ReturnsOneResultPerTest(t *testing.T) {
	ctx := context.NewContext()
	steps := []step.Step{step.NewAlwaysPass(), step.NewAlwaysPass()}
	tc := NewTestCase("test case", steps)

	result := tc.Run(ctx)

	assert.Equal(t, "test case", result.Name)
	assert.Len(t, result.Results, 2)
}
