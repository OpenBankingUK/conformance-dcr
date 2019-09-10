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
					ParseClientRegisterResponse(authoriser).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve client credentials grant").
					WithHttpClient(secureClient).
					GetClientCredentialsGrant(openIDConfig.TokenEndpoint).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Delete software client").
					WithHttpClient(secureClient).
					ClientDelete(openIDConfig.RegistrationEndpointAsString()).
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
					ParseClientRegisterResponse(authoriser).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve client credentials grant").
					WithHttpClient(secureClient).
					GetClientCredentialsGrant(openIDConfig.TokenEndpoint).
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
					Build(),
			).
			Build(),
		NewBuilder("I should not be able to retrieve a registered software if I send invalid credentials").
			TestCase(
				NewTestCaseBuilder("Register software client").
					WithHttpClient(secureClient).
					GenerateSignedClaims(authoriser).
					PostClientRegister(openIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeCreated().
					ParseClientRegisterResponse(authoriser).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve software client with invalid credentials should not succeed").
					WithHttpClient(secureClient).
					ClientRetrieveWithInvalidToken(openIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeUnauthorized().
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve client credentials grant").
					WithHttpClient(secureClient).
					GetClientCredentialsGrant(openIDConfig.TokenEndpoint).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Delete software client").
					WithHttpClient(secureClient).
					ClientDelete(openIDConfig.RegistrationEndpointAsString()).
					Build(),
			).
			Build(),
	}
}
