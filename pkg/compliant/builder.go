package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"net/http"
)

type Builder struct {
	name string
	tcs  []TestCase
}

func NewBuilder(name string) *Builder {
	return &Builder{
		name: name,
		tcs:  []TestCase{},
	}
}

func (b *Builder) TestCase(tc TestCase) *Builder {
	b.tcs = append(b.tcs, tc)
	return b
}

func (b *Builder) Build() Scenario {
	return NewScenario(b.name, b.tcs)
}

type testCaseBuilder struct {
	name  string
	steps []step.Step
}

func NewTestCaseBuilder(name string) *testCaseBuilder {
	return &testCaseBuilder{
		name:  name,
		steps: []step.Step{},
	}
}

func (t *testCaseBuilder) Get(url string) *testCaseBuilder {
	t.steps = append(t.steps, step.NewGetRequest(url, "response", &http.Client{}))
	return t
}

func (t *testCaseBuilder) AssertStatusCodeOk() *testCaseBuilder {
	nextStep := step.NewAssertStatus(200, "response")
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) AssertContextTypeApplicationHtml() *testCaseBuilder {
	nextStep := step.NewAssertContentType("response", "application/html")
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ParseWellKnownRegistrationEndpoint() *testCaseBuilder {
	nextStep := step.NewParseWellKnownRegistrationEndpoint("response", "registration_endpoint")
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) Step(nextStep step.Step) *testCaseBuilder {
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) Build() TestCase {
	return NewTestCase(t.name, t.steps)
}
