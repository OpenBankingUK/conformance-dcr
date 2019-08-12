package compliant

func NewDCR31() Scenarios {
	return Scenarios{
		NewBuilder("Dynamically create a new software client").
			TestCase(
				NewTestCaseBuilder("Creates software client").
					Get("/register").
					AssertStatusCodeOk().
					AssertContextTypeApplicationHtml().
					Build(),
			).
			Build(),
	}
}
