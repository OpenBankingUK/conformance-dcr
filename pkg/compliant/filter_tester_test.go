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
	tester := filteredTester{expression: "Two"}

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

	tester.Compliant(scenarios)

	assert.Equal(t, 1, mockDownstreamTester.count)
}

type mockTester struct {
	count int
}

func (m *mockTester) Name() string { return "mock tester" }

func (m *mockTester) Compliant(scenarios Scenarios) bool {
	m.count = len(scenarios)
	return true
}
