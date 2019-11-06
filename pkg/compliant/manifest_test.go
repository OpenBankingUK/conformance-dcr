package compliant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewManifest(t *testing.T) {
	manifest, err := NewManifest("DCR", "1.0", Scenarios{})

	assert.NoError(t, err)
	expected := versionedManifest{
		name:      "DCR",
		version:   "1.0",
		scenarios: Scenarios{},
	}
	assert.Equal(t, expected, manifest)
}

func TestNewManifestErrorsOnDuplicatedScenarioIds(t *testing.T) {
	scenarios := Scenarios{
		scenario{id: "1"},
		scenario{id: "1"},
	}

	manifest, err := NewManifest("DCR", "1.0", scenarios)

	assert.EqualError(t, err, "scenario must have unique ids")
	assert.Equal(t, versionedManifest{}, manifest)
}

func TestRunReturnsOneResultPerScenario(t *testing.T) {
	scenarios := Scenarios{
		scenario{id: "1"},
		scenario{id: "2"},
	}
	manifest, err := NewManifest("DCR", "1.0", scenarios)
	assert.NoError(t, err)

	result := manifest.Run()

	assert.Len(t, result.Results, 2)
}

func TestRunPassOnNoFails(t *testing.T) {
	scenarios := Scenarios{
		scenario{id: "1"},
	}
	manifest, err := NewManifest("DCR", "1.0", scenarios)
	assert.NoError(t, err)

	result := manifest.Run()

	assert.False(t, result.Fail())
}

func TestNewFilteredManifest_ById(t *testing.T) {
	scenarios := Scenarios{
		scenario{id: "1", name: "name one"},
		scenario{id: "2", name: "name two"},
	}
	manifest, err := NewManifest("DCR", "1.0", scenarios)
	assert.NoError(t, err)

	manifest, err = NewFilteredManifest(manifest, "1")
	assert.NoError(t, err)

	filteredScenarios := manifest.Scenarios()
	assert.Len(t, filteredScenarios, 1)
	assert.Equal(t, "name one", filteredScenarios[0].Name())
}

func TestNewFilteredManifest_ByName(t *testing.T) {
	scenarios := Scenarios{
		scenario{id: "1", name: "name one"},
		scenario{id: "2", name: "name two"},
	}
	manifest, err := NewManifest("DCR", "1.0", scenarios)
	assert.NoError(t, err)

	manifest, err = NewFilteredManifest(manifest, "Two")
	assert.NoError(t, err)

	filteredScenarios := manifest.Scenarios()
	assert.Len(t, filteredScenarios, 1)
	assert.Equal(t, "name two", filteredScenarios[0].Name())
}

func TestNewFilteredManifest_ErrorsOnNoTests(t *testing.T) {
	scenarios := Scenarios{
		scenario{id: "1", name: "name one"},
		scenario{id: "2", name: "name two"},
	}
	manifest, err := NewManifest("DCR", "1.0", scenarios)
	assert.NoError(t, err)

	filteredManifest, err := NewFilteredManifest(manifest, "3")

	assert.EqualError(t, err, "no tests found to run")
	assert.Nil(t, filteredManifest)
}
