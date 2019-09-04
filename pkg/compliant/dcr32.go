package compliant

import (
	"net/http"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
)

func NewDCR32(
	wellKnownEndpoint string,
	openIDConfig openid.Configuration,
	secureClient *http.Client,
	authoriser auth.Authoriser,
) Scenarios {
	return Scenarios{
		NewBuilder("Validate OIDC Config Registration URL").
			TestCase(
				NewTestCaseBuilder("Validate Registration URL").
					ValidateRegistrationEndpoint(openIDConfig.RegistrationEndpoint).
					Build(),
			).
			Build(),
		NewBuilder("Dynamically create a new software client").
			TestCase(
				NewTestCaseBuilder("Register software client").
					WithHttpClient(secureClient).
					GenerateSignedClaims(authoriser).
					PostClientRegister(openIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeCreated().
					ParseClientRegisterResponse().
					Build(),
			).
			Build(),
		NewBuilder("Dynamically retrieve a new software client").
			TestCase(
				NewTestCaseBuilder("Register software client").
					WithHttpClient(secureClient).
					GenerateSignedClaims(authoriser).
					PostClientRegister(openIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeCreated().
					ParseClientRegisterResponse().
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve software client").
					WithHttpClient(secureClient).
					ClientRetrieve(openIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeOk().
					ParseClientRetrieveResponse().
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Delete software client").
					WithHttpClient(secureClient).
					ClientDelete(openIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeOk().
					Build(),
			).
			Build(),
	}
}
