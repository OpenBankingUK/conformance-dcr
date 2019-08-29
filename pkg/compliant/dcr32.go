package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"net/http"
)

func NewDCR32(
	wellKnownEndpoint, registrationEndpoint string,
	secureClient *http.Client,
	authoriser auth.Authoriser,
) Scenarios {
	return Scenarios{
		NewBuilder("Dynamically create a new software client").
			TestCase(
				NewTestCaseBuilder("Register software client").
					WithHttpClient(secureClient).
					GenerateSignedClaims(authoriser).
					PostClientRegister(registrationEndpoint).
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
					PostClientRegister(registrationEndpoint).
					AssertStatusCodeCreated().
					ParseClientRegisterResponse().
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve software client").
					WithHttpClient(secureClient).
					ClientRetrieve(registrationEndpoint).
					AssertStatusCodeOk().
					ParseClientRetrieveResponse().
					Build(),
			).
			Build(),
	}
}
