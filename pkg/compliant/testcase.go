package compliant

import (
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/step"
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
	results := make(step.Results, 0, len(t.steps))
	for _, step := range t.steps {
		results = append(results, step.Run(ctx))
	}

	return TestCaseResult{
		Name:    t.name,
		Results: results,
	}
}
