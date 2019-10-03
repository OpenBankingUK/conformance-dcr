package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScenarioResult(t *testing.T) {
	results := ScenariosResult{
		ScenarioResult{
			Name: "some scenario name",
			TestCaseResults: TestCaseResults{
				TestCaseResult{
					Name: "some test case",
					Results: step.Results{
						step.Result{
							Name:       "some step",
							Pass:       true,
							FailReason: "",
						},
					},
				},
			},
		},
	}

	assert.False(t, results.Fail())
}

func TestScenarioResult_OneFails(t *testing.T) {
	results := ScenariosResult{
		ScenarioResult{
			Name: "some scenario name",
			TestCaseResults: TestCaseResults{
				TestCaseResult{
					Name: "some test case",
					Results: step.Results{
						step.Result{
							Name:       "some step that fails",
							Pass:       false,
							FailReason: "",
						},
						step.Result{
							Name:       "some step that passes",
							Pass:       true,
							FailReason: "",
						},
					},
				},
			},
		},
	}

	assert.True(t, results.Fail())
}

func TestNewScenario_RunsAllTestCases(t *testing.T) {
	tcs := []TestCase{
		testCase{
			name:  "test case 1",
			steps: []step.Step{},
		},
		testCase{
			name:  "test case 2",
			steps: []step.Step{},
		},
	}
	scenario := NewScenario("some scenario", "spec link", tcs)

	results := scenario.Run()

	assert.Equal(t, "some scenario", scenario.Name())
	assert.Equal(t, "spec link", scenario.Spec())
	assert.False(t, results.Fail())
	assert.Equal(t, "some scenario", results.Name)
	assert.Len(t, results.TestCaseResults, 2)
}
