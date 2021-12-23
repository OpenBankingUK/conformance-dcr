package compliant

import (
	"errors"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/step"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerboseTester_Compliant(t *testing.T) {
	scenarios := Scenarios{
		NewBuilder("#1", "Scenario with one test", "Spec Link").
			TestCase(
				NewTestCaseBuilder("Always fail test").
					Step(failStep{}).
					Build(),
			).
			Build(),
	}
	manifest, err := NewManifest("test", "0.0", scenarios)
	assert.NoError(t, err)
	tester := NewTester()

	isCompliant, err := tester.Compliant(manifest)

	assert.NoError(t, err)
	assert.False(t, isCompliant)
}

type failStep struct{}

func (s failStep) Run(ctx step.Context) step.Result {
	return step.NewFailResult("test name", "reasons")
}

func TestVerboseTester_CallsListener(t *testing.T) {
	scenarios := Scenarios{
		NewBuilder("#1", "Scenario with one test", "Spec Link").
			TestCase(
				NewTestCaseBuilder("Always pass test").
					Step(passStep{}).
					Build(),
			).
			Build(),
	}
	manifest, err := NewManifest("test", "0.0", scenarios)
	assert.NoError(t, err)
	tester := NewTester()
	called := false
	tester.AddListener(func(result ManifestResult) error {
		called = true
		return nil
	})

	isCompliant, err := tester.Compliant(manifest)

	assert.NoError(t, err)
	assert.True(t, isCompliant)
	assert.True(t, called)
}

func TestVerboseTester_HandlesListenerError(t *testing.T) {
	scenarios := Scenarios{
		NewBuilder("#1", "Scenario with one test", "Spec Link").
			TestCase(
				NewTestCaseBuilder("Always pass test").
					Step(passStep{}).
					Build(),
			).
			Build(),
	}
	manifest, err := NewManifest("test", "0.0", scenarios)
	assert.NoError(t, err)
	tester := NewTester()
	tester.AddListener(func(result ManifestResult) error {
		return errors.New("boom")
	})

	isCompliant, err := tester.Compliant(manifest)

	assert.EqualError(t, err, "boom")
	assert.False(t, isCompliant)
}
