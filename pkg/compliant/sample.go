package compliant

import "bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"

func NewGoogleScenarios() []Scenario {
	return []Scenario{
		NewScenario(
			"Google works",
			[]TestCase{
				NewTestCase(
					1,
					"Google landing page is reachable",
					[]step.Step{
						step.NewGetRequest(1, "https://www.google.com", "response"),
						step.NewAssertStatusOk(2, "response"),
						step.NewAssertContentType(3, "response", "application/html"),
					},
				),
			},
		),
	}
}
