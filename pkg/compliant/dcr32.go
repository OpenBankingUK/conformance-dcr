package compliant

func NewDCR32(wellKnownEndpoint, ssa string) Scenarios {
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
					ClientRegister(ssa).
					AssertStatusCodeCreated().
					ParseClientRegisterResponse().
					Build(),
			).
			Build(),
	}
}
