package compliant

import (
	"crypto/rsa"
	"net/http"
)

func NewDCR32(wellKnownEndpoint, ssa string, privateKey *rsa.PrivateKey, secureClient *http.Client) Scenarios {
	return Scenarios{
		NewBuilder("Dynamically create a new software client").
			TestCase(
				NewTestCaseBuilder("Retrieve registration endpoint from OIDC Discovery Endpoint").
					Get(wellKnownEndpoint).
					AssertStatusCodeOk().
					ParseWellKnownRegistrationEndpoint().
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Register software client").
					WithHttpClient(secureClient).
					GenerateSignedClaims(ssa, privateKey).
					PostClientRegister().
					AssertStatusCodeCreated().
					ParseClientRegisterResponse().
					Build(),
			).
			Build(),
		NewBuilder("Dynamically retrieve a new software client").
			TestCase(
				NewTestCaseBuilder("Retrieve registration endpoint from OIDC Discovery Endpoint").
					Get(wellKnownEndpoint).
					AssertStatusCodeOk().
					ParseWellKnownRegistrationEndpoint().
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Register software client").
					WithHttpClient(secureClient).
					GenerateSignedClaims(ssa, privateKey).
					PostClientRegister().
					AssertStatusCodeCreated().
					ParseClientRegisterResponse().
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve software client").
					WithHttpClient(secureClient).
					ClientRetrieve().
					AssertStatusCodeOk().
					ParseClientRetrieveResponse().
					Build(),
			).
			Build(),
	}
}
