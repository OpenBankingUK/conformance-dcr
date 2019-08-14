package compliant

func NewDCR31(wellKnownEndpoint string) Scenarios {
	return Scenarios{
		NewBuilder("Dynamically create a new software client").
			TestCase(
				NewTestCaseBuilder("Retrieve registration endpoint from OIDC Discovery Endpoint").
					Get("https://modelobankauth2018.o3bank.co.uk:4101/.well-known/openid-configuration").
					AssertStatusCodeOk().
					ParseWellKnownRegistrationEndpoint().
					Build(),
			).
			Build(),
	}
}
