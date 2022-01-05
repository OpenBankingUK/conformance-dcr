package compliant

import (
	"testing"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/step"
	"github.com/stretchr/testify/assert"
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
	ctx := step.NewContext()
	steps := []step.Step{passStep{}, passStep{}}
	tc := NewTestCase("test case", steps)

	result := tc.Run(ctx)

	assert.Equal(t, "test case", result.Name)
	assert.Len(t, result.Results, 2)
}

type passStep struct{}

func (s passStep) Run(ctx step.Context) step.Result {
	return step.NewPassResult("test name")
}
