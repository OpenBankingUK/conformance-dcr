package compliant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilteredTester_FilterTheRightScenario(t *testing.T) {
	scenarios := Scenarios{
		NewScenario("#1", "One", "spec link", nil),
		NewScenario("#2", "Two", "spec link", nil),
	}
	tester := filteredTester{expression: "two"}

	filteredScenarios := tester.filter(scenarios)

	assert.Len(t, filteredScenarios, 1)
	assert.Equal(t, "Two", filteredScenarios[0].Name())
}

func TestFilteredTester_ShouldRunOnlyOneScenario(t *testing.T) {
	scenarios := Scenarios{
		NewBuilder("#1", "Scenario with ONE test", "Spec link").
			TestCase(
				NewTestCaseBuilder("Always fail test").
					Step(failStep{}).
					Build(),
			).
			Build(),
		NewBuilder("#2", "Scenario with TWO test", "Spec link").
			TestCase(
				NewTestCaseBuilder("Always fail test").
					Step(failStep{}).
					Build(),
			).
			Build(),
	}
	manifest, err := NewManifest("test", "0.0", scenarios)
	assert.NoError(t, err)

	mockDownstreamTester := &mockTester{count: 1}
	tester := NewFilteredTester("TWO", mockDownstreamTester)

	_, err = tester.Compliant(manifest)

	assert.NoError(t, err)
	assert.Equal(t, 1, mockDownstreamTester.count)
}

type mockTester struct {
	count int
}

func (m *mockTester) Name() string { return "mock tester" }

func (m *mockTester) Compliant(manifest Manifest) (bool, error) {
	m.count = len(manifest.scenarios)
	return true, nil
}

func TestFilteredTester_ReturnsErrorIfNoTestsFilterTheRightScenario(t *testing.T) {
	mockDownstreamTester := &mockTester{count: 0}
	tester := NewFilteredTester("TWO", mockDownstreamTester)

	_, err := tester.Compliant(Manifest{})

	assert.EqualError(t, err, "no tests found to run")
}
