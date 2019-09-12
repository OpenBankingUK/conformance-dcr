package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilteredTester_FilterTheRightScenario(t *testing.T) {
	scenarios := Scenarios{
		NewScenario("One", nil),
		NewScenario("Two", nil),
	}
	tester := filteredTester{expression: "two"}

	filteredScenarios := tester.filter(scenarios)

	assert.Len(t, filteredScenarios, 1)
	assert.Equal(t, "Two", filteredScenarios[0].Name())
}

func TestFilteredTester_ShouldRunOnlyOneScenario(t *testing.T) {
	scenarios := Scenarios{
		NewBuilder("Scenario with ONE test").
			TestCase(
				NewTestCaseBuilder("Always fail test").
					Step(step.NewAlwaysFail()).
					Build(),
			).
			Build(),
		NewBuilder("Scenario with TWO test").
			TestCase(
				NewTestCaseBuilder("Always fail test").
					Step(step.NewAlwaysFail()).
					Build(),
			).
			Build(),
	}
	mockDownstreamTester := &mockTester{count: 1}
	tester := NewFilteredTester("TWO", mockDownstreamTester)

	_, err := tester.Compliant(scenarios)

	assert.NoError(t, err)
	assert.Equal(t, 1, mockDownstreamTester.count)
}

type mockTester struct {
	count int
}

func (m *mockTester) Name() string { return "mock tester" }

func (m *mockTester) Compliant(scenarios Scenarios) (bool, error) {
	m.count = len(scenarios)
	return true, nil
}

func TestFilteredTester_ReturnsErrorIfNoTestsFilterTheRightScenario(t *testing.T) {
	scenarios := Scenarios{}

	mockDownstreamTester := &mockTester{count: 0}
	tester := NewFilteredTester("TWO", mockDownstreamTester)

	_, err := tester.Compliant(scenarios)

	assert.EqualError(t, err, "no tests found to run")
}
