package compliant

import (
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/schema"
	"net/http"
	"time"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/step"
)

type Builder struct {
	id   string
	name string
	spec string
	tcs  []TestCase
}

func NewBuilder(id, name, spec string) *Builder {
	return &Builder{
		id:   id,
		name: name,
		spec: spec,
		tcs:  []TestCase{},
	}
}

func (b *Builder) TestCase(tc ...TestCase) *Builder {
	b.tcs = append(b.tcs, tc...)
	return b
}

func (b *Builder) Build() Scenario {
	return NewScenario(b.id, b.name, b.spec, b.tcs)
}

type testCaseBuilder struct {
	name       string
	steps      []step.Step
	httpClient *http.Client
}

func NewTestCaseBuilder(name string) *testCaseBuilder {
	return &testCaseBuilder{
		name:       name,
		steps:      []step.Step{},
		httpClient: newDefaultHttpClient(),
	}
}

func newDefaultHttpClient() *http.Client {
	return &http.Client{Timeout: time.Second * 10}
}

const (
	responseCtxKey   = "response"
	clientCtxKey     = "software_client"
	jwtClaimsCtxKey  = "jwt_claims"
	grantTokenCtxKey = "grant_token"
)

func (t *testCaseBuilder) WithHttpClient(client *http.Client) *testCaseBuilder {
	t.httpClient = client
	return t
}

func (t *testCaseBuilder) Get(url string) *testCaseBuilder {
	t.steps = append(t.steps, step.NewGetRequest(url, responseCtxKey, t.httpClient))
	return t
}

func (t *testCaseBuilder) AssertStatusCodeOk() *testCaseBuilder {
	nextStep := step.NewAssertStatus(http.StatusOK, responseCtxKey)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) AssertStatusCodeUnauthorized() *testCaseBuilder {
	nextStep := step.NewAssertStatus(http.StatusUnauthorized, responseCtxKey)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) AssertStatusCodeBadRequest() *testCaseBuilder {
	nextStep := step.NewAssertStatus(http.StatusBadRequest, responseCtxKey)
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

func (t *testCaseBuilder) GenerateSignedClaims(authoriserBuilder auth.AuthoriserBuilder) *testCaseBuilder {
	nextStep := step.NewClaims(jwtClaimsCtxKey, authoriserBuilder)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) PostClientRegister(registrationEndpoint string) *testCaseBuilder {
	nextStep := step.NewPostClientRegister(registrationEndpoint, jwtClaimsCtxKey, responseCtxKey, t.httpClient)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ClientUpdate(registrationEndpoint string) *testCaseBuilder {
	nextStep := step.NewClientUpdate(
		registrationEndpoint,
		jwtClaimsCtxKey,
		responseCtxKey,
		clientCtxKey,
		grantTokenCtxKey,
		t.httpClient,
	)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ClientDelete(registrationEndpoint string) *testCaseBuilder {
	nextStep := step.NewClientDelete(registrationEndpoint, clientCtxKey, grantTokenCtxKey, t.httpClient)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ClientRetrieve(registrationEndpoint string) *testCaseBuilder {
	nextStep := step.NewClientRetrieve(responseCtxKey, registrationEndpoint, clientCtxKey, grantTokenCtxKey, t.httpClient)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) AssertValidSchemaResponse(validator schema.Validator) *testCaseBuilder {
	nextStep := step.NewClientRetrieveSchema(responseCtxKey, validator)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) SetInvalidGrantToken() *testCaseBuilder {
	nextStep := step.NewSetInvalidGrantToken(grantTokenCtxKey)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ParseClientRegisterResponse(authoriserBuilder auth.AuthoriserBuilder) *testCaseBuilder {
	nextStep := step.NewClientRegisterResponse(responseCtxKey, clientCtxKey, authoriserBuilder)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ParseClientRetrieveResponse(openIDConfigTokenEndpoint string) *testCaseBuilder {
	nextStep := step.NewClientRetrieveResponse(responseCtxKey, clientCtxKey, openIDConfigTokenEndpoint)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) ValidateRegistrationEndpoint(registrationEndpoint *string) *testCaseBuilder {
	nextStep := step.NewValidateRegistrationEndpoint(registrationEndpoint)
	t.steps = append(t.steps, nextStep)
	return t
}

func (t *testCaseBuilder) GetClientCredentialsGrant(tokenEndpoint string) *testCaseBuilder {
	nextStep := step.NewClientCredentialsGrant(grantTokenCtxKey, clientCtxKey, tokenEndpoint, t.httpClient)
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
