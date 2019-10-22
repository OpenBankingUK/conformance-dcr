package compliant

import "errors"

type Manifest struct {
	name      string
	version   string
	scenarios Scenarios
}

func NewManifest(name, version string, scenarios Scenarios) (Manifest, error) {
	for _, scenario := range scenarios {
		if scenarioIdDuplicated(scenario.Id(), scenarios) {
			return Manifest{}, errors.New("scenario must have unique ids")
		}
	}

	return Manifest{
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

func (s Manifest) Run() ManifestResult {
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
