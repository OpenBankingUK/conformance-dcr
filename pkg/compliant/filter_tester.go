package compliant

import "strings"

type filteredTester struct {
	expression       string
	downstreamTester Tester
}

func NewFilteredTester(expression string, downstreamTester Tester) Tester {
	return filteredTester{
		expression:       expression,
		downstreamTester: downstreamTester,
	}
}

func (f filteredTester) Compliant(scenarios Scenarios) bool {
	filteredScenarios := f.filter(scenarios)
	return f.downstreamTester.Compliant(filteredScenarios)
}

func (f filteredTester) filter(scenarios Scenarios) Scenarios {
	var filteredScenarios Scenarios
	for _, scenario := range scenarios {
		if strings.Contains(scenario.Name(), f.expression) {
			filteredScenarios = append(filteredScenarios, scenario)
		}
	}
	return filteredScenarios
}
