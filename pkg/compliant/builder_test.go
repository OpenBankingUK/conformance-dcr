package compliant

import (
	"crypto/rsa"
	"net/http"
	"testing"
	"time"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/schema"
	"github.com/stretchr/testify/require"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/openid"
	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	scenario := NewBuilder("#1", "scenario name", "spec link").
		TestCase(NewTestCase("some test", nil)).
		TestCase(NewTestCase("another test", nil))

	assert.Equal(t, "scenario name", scenario.name)
	assert.Equal(t, "spec link", scenario.spec)
	assert.Len(t, scenario.tcs, 2)
}

func TestNewTestCaseBuilder(t *testing.T) {
	authoriserBuilder := auth.NewAuthoriserBuilder().
		WithIssuer("issuer").
		WithKID("kid").
		WithSSA("ssa").
		WithPrivateKey(&rsa.PrivateKey{}).
		WithOpenIDConfig(openid.Configuration{}).
		WithJwtExpiration(time.Hour)

	validator, err := schema.NewValidator("3.2")
	require.NoError(t, err)

	const sampleEndpoint = "http://host/path"
	var someUrl *string
	tc := NewTestCaseBuilder("test case").
		WithHttpClient(&http.Client{}).
		Get("www.google.com").
		AssertStatusCodeOk().
		AssertStatusCodeUnauthorized().
		AssertStatusCodeBadRequest().
		AssertStatusCodeCreated().
		AssertContextTypeApplicationHtml().
		GenerateSignedClaims(authoriserBuilder).
		PostClientRegister(sampleEndpoint).
		ParseClientRegisterResponse(authoriserBuilder).
		ClientRetrieve(sampleEndpoint).
		ClientDelete(sampleEndpoint).
		ParseClientRetrieveResponse(sampleEndpoint).
		AssertValidSchemaResponse(validator).
		SetInvalidGrantToken().
		ValidateRegistrationEndpoint(someUrl).
		GetClientCredentialsGrant(sampleEndpoint)

	assert.Equal(t, "test case", tc.name)
	assert.Len(t, tc.steps, 16)
}
