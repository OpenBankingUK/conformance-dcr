package compliant

import (
	"crypto/rsa"
	"net/http"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
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
	jwtClaimsCtxKey    = "jwt_claims"
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

func (t *testCaseBuilder) GenerateSignedClaims(ssa string, privateKey *rsa.PrivateKey) *testCaseBuilder {
	nextStep := step.NewClaims(jwtClaimsCtxKey, openIdConfigCtxKey, ssa, privateKey)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ClientRegister() *testCaseBuilder {
	nextStep := step.NewClientRegister(openIdConfigCtxKey, jwtClaimsCtxKey, responseCtxKey, &http.Client{})
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ClientRetrieve() *testCaseBuilder {
	nextStep := step.NewClientRetrieve(responseCtxKey, openIdConfigCtxKey, clientCtxKey, &http.Client{})
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ParseClientRegisterResponse() *testCaseBuilder {
	nextStep := step.NewClientRegisterResponse(responseCtxKey, clientCtxKey)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ParseClientRetrieveResponse() *testCaseBuilder {
	nextStep := step.NewClientRetrieveResponse(responseCtxKey, clientCtxKey)
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
