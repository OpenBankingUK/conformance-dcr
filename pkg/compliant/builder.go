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

const (
	openIdConfigCtxKey = "openid_config"
	responseCtxKey     = "response"
	clientCtxKey       = "software_client"
)

func (t *testCaseBuilder) Get(url string) *testCaseBuilder {
	t.steps = append(t.steps, step.NewGetRequest(url, responseCtxKey, &http.Client{}))
	return t
}

func (t *testCaseBuilder) AssertStatusCodeOk() *testCaseBuilder {
	nextStep := step.NewAssertStatus(http.StatusOK, responseCtxKey)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) AssertStatusCodeCreated() *testCaseBuilder {
	nextStep := step.NewAssertStatus(http.StatusCreated, responseCtxKey)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) AssertContextTypeApplicationHtml() *testCaseBuilder {
	nextStep := step.NewAssertContentType(responseCtxKey, "application/html")
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ParseWellKnownRegistrationEndpoint() *testCaseBuilder {
	nextStep := step.NewParseWellKnownRegistrationEndpoint(responseCtxKey, openIdConfigCtxKey)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ClientRegister(ssa string) *testCaseBuilder {
	nextStep := step.NewClientRegister(openIdConfigCtxKey, ssa, responseCtxKey, &http.Client{})
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ParseClientRegisterResponse() *testCaseBuilder {
	nextStep := step.NewClientRegisterResponse(responseCtxKey, clientCtxKey)
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
