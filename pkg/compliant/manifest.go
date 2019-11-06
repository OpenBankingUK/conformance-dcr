package compliant

import (
	"errors"
	"fmt"
	"strings"
)

type Manifest interface {
	Run() ManifestResult
	Scenarios() Scenarios
	Name() string
	Version() string
}

type versionedManifest struct {
	name      string
	version   string
	scenarios Scenarios
}

func NewManifest(name, version string, scenarios Scenarios) (Manifest, error) {
	for _, scenario := range scenarios {
		if scenarioIdDuplicated(scenario.Id(), scenarios) {
			return versionedManifest{}, errors.New("scenario must have unique ids")
		}
	}

	return versionedManifest{
		name:      name,
		version:   version,
		scenarios: scenarios,
	}, nil
}

func scenarioIdDuplicated(id string, scenarios Scenarios) bool {
	count := 0
	for _, scenario := range scenarios {
		if scenario.Id() == id {
			count++
		}
	}
	return count != 1
}

func (s versionedManifest) Run() ManifestResult {
	results := make([]ScenarioResult, len(s.scenarios))
	for key, scenario := range s.scenarios {
		results[key] = scenario.Run()
	}
	return ManifestResult{
		Results: results,
		Name:    s.name,
		Version: s.version,
	}
}

func (s versionedManifest) Scenarios() Scenarios {
	return s.scenarios
}

func (s versionedManifest) Name() string {
	return s.name
}

func (s versionedManifest) Version() string {
	return s.version
}

type ManifestResult struct {
	Results []ScenarioResult
	Name    string
	Version string
}

func (r ManifestResult) Fail() bool {
	for _, result := range r.Results {
		if result.Fail() {
			return true
		}
	}
	return false
}

func NewFilteredManifest(manifest Manifest, expression string) (Manifest, error) {
	scenarios := manifest.Scenarios()

	filteredScenarios := filter(scenarios, expression)
	if len(filteredScenarios) == 0 {
		return nil, errors.New("no tests found to run")
	}

	return NewManifest(
		fmt.Sprintf("(filtered) %s", manifest.Name()),
		manifest.Version(),
		filteredScenarios,
	)
}

func filter(scenarios Scenarios, expression string) Scenarios {
	var filteredScenarios Scenarios
	for _, scenario := range scenarios {
		if scenarioNameContains(scenario, expression) ||
			scenarioIdContains(scenario, expression) {
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
