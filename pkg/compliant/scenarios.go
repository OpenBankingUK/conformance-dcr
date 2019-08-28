package compliant

import "bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"

type Scenario interface {
	Run() ScenarioResult
	Name() string
}

type Scenarios []Scenario

type ScenarioResult struct {
	Name string
	TestCaseResults
}

type ScenariosResult []ScenarioResult

func (r ScenariosResult) Fail() bool {
	for _, result := range r {
		if result.Fail() {
			return true
		}
	}
	return false
}

type scenario struct {
	name string
	tcs  []TestCase
}

func NewScenario(name string, tcs []TestCase) Scenario {
	return scenario{
		name: name,
		tcs:  tcs,
	}
}

func (s scenario) Name() string {
	return s.name
}

func (s scenario) Run() ScenarioResult {
	ctx := step.NewContext()
	var results TestCaseResults
	for _, tc := range s.tcs {
		tcResult := tc.Run(ctx)
		results = append(results, tcResult)
	}

	return ScenarioResult{
		Name:            s.name,
		TestCaseResults: results,
	}
}
