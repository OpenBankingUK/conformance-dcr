package compliant

func NewDCR31(wellKnownEndpoint string) Scenarios {
	return Scenarios{
		NewBuilder("Dynamically create a new software client").
			TestCase(
				NewTestCaseBuilder("Retrieve registration endpoint from OIDC Discovery Endpoint").
					Get(wellKnownEndpoint).
					AssertStatusCodeOk().
					ParseWellKnownRegistrationEndpoint().
					Build(),
			).
			Build(),
	}
}
