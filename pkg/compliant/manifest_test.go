package compliant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewManifest(t *testing.T) {
	manifest, err := NewManifest("DCR", "1.0", Scenarios{})

	assert.NoError(t, err)
	expected := Manifest{
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
	assert.Equal(t, Manifest{}, manifest)
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
