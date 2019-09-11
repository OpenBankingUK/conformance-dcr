package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	scenario := NewBuilder("scenario name").
		TestCase(NewTestCase("some test", nil)).
		TestCase(NewTestCase("another test", nil))

	assert.Equal(t, "scenario name", scenario.name)
	assert.Len(t, scenario.tcs, 2)
}

func TestNewTestCaseBuilder(t *testing.T) {
	authoriser := auth.NewAuthoriser(
		openid.Configuration{},
		"ssa",
		"kid",
		"clientId",
		[]string{},
		&rsa.PrivateKey{},
	)

	const sampleEndpoint = "http://host/path"
	tc := NewTestCaseBuilder("test case").
		WithHttpClient(&http.Client{}).
		Get("www.google.com").
		AssertStatusCodeOk().
		AssertContextTypeApplicationHtml().
		GenerateSignedClaims(authoriser).
		PostClientRegister(sampleEndpoint).
		AssertStatusCodeCreated().
		ParseClientRegisterResponse(authoriser).
		ClientRetrieve(sampleEndpoint).
		ClientDelete(sampleEndpoint).
		ParseClientRetrieveResponse(sampleEndpoint).
		Step(step.NewAlwaysPass())

	assert.Equal(t, "test case", tc.name)
	assert.Len(t, tc.steps, 11)
}
