package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"sort"
)

type TestCase interface {
	Run(ctx context.Context) TestCaseResults
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
	order int
	steps []step.Step
}

func NewTestCase(order int, name string, steps []step.Step) testCase {
	return testCase{
		name:  name,
		order: order,
		steps: steps,
	}
}

func (t testCase) Run(ctx context.Context) TestCaseResults {
	sort.Slice(t.steps, func(i, j int) bool {
		return t.steps[i].Order() < t.steps[j].Order()
	})

	var results step.Results
	for _, step := range t.steps {
		results = append(results, step.Run(ctx))
	}
	return TestCaseResults{
		{
			Name:    t.name,
			Results: results,
		},
	}
}
