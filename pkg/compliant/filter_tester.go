package compliant

import (
	"errors"
	"fmt"
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

func (f filteredTester) Compliant(manifest Manifest) (bool, error) {
	filteredScenarios := f.filter(manifest.scenarios)
	if len(filteredScenarios) == 0 {
		return false, errors.New("no tests found to run")
	}

	filteredManifest, err := NewManifest(
		fmt.Sprintf("(filtered) %s", manifest.name),
		manifest.version,
		filteredScenarios,
	)
	if err != nil {
		return false, err
	}

	return f.downstreamTester.Compliant(filteredManifest)
}

func (f filteredTester) filter(scenarios Scenarios) Scenarios {
	var filteredScenarios Scenarios
	for _, scenario := range scenarios {
		if scenarioNameContains(scenario, f.expression) ||
			scenarioIdContains(scenario, f.expression) {
			filteredScenarios = append(filteredScenarios, scenario)
		}
	}
	return filteredScenarios
}

func scenarioNameContains(scenario Scenario, expression string) bool {
	return strings.Contains(
		strings.ToLower(scenario.Name()),
		strings.ToLower(expression),
	)
}

func scenarioIdContains(scenario Scenario, expression string) bool {
	return strings.Contains(
		strings.ToLower(scenario.Id()),
		strings.ToLower(expression),
	)
}
