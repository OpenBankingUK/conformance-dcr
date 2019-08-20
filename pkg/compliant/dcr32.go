package compliant

import "crypto/rsa"

func NewDCR32(wellKnownEndpoint, ssa string, privateKey *rsa.PrivateKey) Scenarios {
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
					GenerateSignedClaims(ssa, privateKey).
					ClientRegister().
					AssertStatusCodeCreated().
					ParseClientRegisterResponse().
					Build(),
			).
			Build(),
	}
}
