package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
)

type TestCase interface {
	Run(ctx step.Context) TestCaseResult
}

type TestCaseResult struct {
	Name string
	step.Results
}

type TestCaseResults []TestCaseResult

func (r TestCaseResults) Fail() bool {
	for _, result := range r {
		if result.Fail() {
			return true
		}
	}
	return false
}

type testCase struct {
	name  string
	steps []step.Step
}

func NewTestCase(name string, steps []step.Step) testCase {
	return testCase{
		name:  name,
		steps: steps,
	}
}

func (t testCase) Run(ctx step.Context) TestCaseResult {
	var results step.Results
	for _, step := range t.steps {
		results = append(results, step.Run(ctx))
	}

	return TestCaseResult{
		Name:    t.name,
		Results: results,
	}
}
