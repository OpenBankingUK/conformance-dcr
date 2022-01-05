package compliant

import "bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"

type Scenario interface {
	Run() ScenarioResult
	Id() string
	Name() string
	Spec() string
}

type Scenarios []Scenario

type ScenarioResult struct {
	Id   string
	Name string
	Spec string
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
	id   string
	name string
	spec string
	tcs  []TestCase
}

func NewScenario(id, name, spec string, tcs []TestCase) Scenario {
	return scenario{
		id:   id,
		name: name,
		spec: spec,
		tcs:  tcs,
	}
}

func (s scenario) Id() string {
	return s.id
}

func (s scenario) Name() string {
	return s.name
}

func (s scenario) Spec() string {
	return s.spec
}

func (s scenario) Run() ScenarioResult {
	ctx := step.NewContext()
	var results TestCaseResults
	for _, tc := range s.tcs {
		tcResult := tc.Run(ctx)
		results = append(results, tcResult)
	}

	return ScenarioResult{
		Id:              s.id,
		Name:            s.name,
		Spec:            s.spec,
		TestCaseResults: results,
	}
}
