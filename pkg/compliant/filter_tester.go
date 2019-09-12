package compliant

import (
	"errors"
	"strings"
)

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

func (f filteredTester) Compliant(scenarios Scenarios) (bool, error) {
	filteredScenarios := f.filter(scenarios)
	if len(filteredScenarios) == 0 {
		return false, errors.New("no tests found to run")
	}
	return f.downstreamTester.Compliant(filteredScenarios)
}

func (f filteredTester) filter(scenarios Scenarios) Scenarios {
	var filteredScenarios Scenarios
	for _, scenario := range scenarios {
		if strings.Contains(
			strings.ToLower(scenario.Name()),
			strings.ToLower(f.expression),
		) {
			filteredScenarios = append(filteredScenarios, scenario)
		}
	}
	return filteredScenarios
}
