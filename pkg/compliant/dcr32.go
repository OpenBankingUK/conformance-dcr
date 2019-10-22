package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/schema"
	"net/http"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
)

func NewDCR32(
	cfg DCR32Config,
	secureClient *http.Client,
	authoriserBuilder auth.AuthoriserBuilder,
	validator schema.Validator,
) (Manifest, error) {
	// nolint:lll
	const (
		specLinkDiscovery        = "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2#DynamicClientRegistration-v3.2-Discovery"
		specLinkRegisterSoftware = "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2#DynamicClientRegistration-v3.2-POST/register"
		specLinkDeleteSoftware   = "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2#DynamicClientRegistration-v3.2-DELETE/register/{ClientId}"
		specLinkRetrieveSoftware = "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2#DynamicClientRegistration-v3.2-GET/register/{ClientId}"
	)
	scenarios := Scenarios{
		NewBuilder(
			"DCR-001",
			"Validate OIDC Config Registration URL",
			specLinkDiscovery,
		).TestCase(
			NewTestCaseBuilder("Validate Registration URL").
				ValidateRegistrationEndpoint(cfg.OpenIDConfig.RegistrationEndpoint).
				Build(),
		).
			Build(),
		NewBuilder(
			"DCR-002",
			"Dynamically create a new software client",
			specLinkRegisterSoftware,
		).
			TestCase(
				NewTestCaseBuilder("Register software client").
					WithHttpClient(secureClient).
					GenerateSignedClaims(authoriserBuilder).
					PostClientRegister(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeCreated().
					ParseClientRegisterResponse(authoriserBuilder).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve client credentials grant").
					WithHttpClient(secureClient).
					GetClientCredentialsGrant(cfg.OpenIDConfig.TokenEndpoint).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Delete software client").
					WithHttpClient(secureClient).
					ClientDelete(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					Build(),
			).
			Build(),
		NewBuilder(
			"DCR-003",
			"Delete software statement is supported",
			specLinkDeleteSoftware,
		).
			TestCase(
				NewTestCaseBuilder("Register software client").
					WithHttpClient(secureClient).
					GenerateSignedClaims(authoriserBuilder).
					PostClientRegister(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeCreated().
					ParseClientRegisterResponse(authoriserBuilder).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve client credentials grant").
					WithHttpClient(secureClient).
					GetClientCredentialsGrant(cfg.OpenIDConfig.TokenEndpoint).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Delete software client").
					WithHttpClient(secureClient).
					ClientDelete(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve delete software client should fail").
					WithHttpClient(secureClient).
					ClientRetrieve(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeUnauthorized().
					Build(),
			).
			Build(),
		NewBuilder(
			"DCR-004",
			"Dynamically create a new software client will fail on invalid registration request",
			specLinkRegisterSoftware,
		).
			TestCase(
				NewTestCaseBuilder("Register software client fails on expired claims").
					WithHttpClient(secureClient).
					GenerateSignedClaims(
						authoriserBuilder.
							WithJwtExpiration(-time.Hour),
					).
					PostClientRegister(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeBadRequest().
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Register software client fails on invalid issuer").
					WithHttpClient(secureClient).
					GenerateSignedClaims(
						authoriserBuilder.
							WithOpenIDConfig(
								openid.NewBuilder().
									From(cfg.OpenIDConfig).
									WithIssuer("foo.is/invalid").
									Build(),
							),
					).
					PostClientRegister(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeBadRequest().
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Register software client fails on invalid issuer too short").
					WithHttpClient(secureClient).
					GenerateSignedClaims(
						authoriserBuilder.
							WithOpenIDConfig(
								openid.NewBuilder().
									From(cfg.OpenIDConfig).
									WithIssuer("").
									Build(),
							),
					).
					PostClientRegister(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeBadRequest().
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Register software client fails on invalid issuer too long").
					WithHttpClient(secureClient).
					GenerateSignedClaims(
						authoriserBuilder.
							WithOpenIDConfig(
								openid.NewBuilder().
									From(cfg.OpenIDConfig).
									WithIssuer("123456789012345678901234567890").
									Build(),
							),
					).
					PostClientRegister(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeBadRequest().
					Build(),
			).
			Build(),
		NewBuilder(
			"DCR-005",
			"Dynamically retrieve a new software client",
			specLinkRetrieveSoftware,
		).
			TestCase(
				NewTestCaseBuilder("Register software client").
					WithHttpClient(secureClient).
					GenerateSignedClaims(authoriserBuilder).
					PostClientRegister(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeCreated().
					ParseClientRegisterResponse(authoriserBuilder).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve client credentials grant").
					WithHttpClient(secureClient).
					GetClientCredentialsGrant(cfg.OpenIDConfig.TokenEndpoint).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve software client").
					WithHttpClient(secureClient).
					ClientRetrieve(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeOk().
					AssertValidSchemaResponse(validator).
					ParseClientRetrieveResponse(cfg.OpenIDConfig.TokenEndpoint).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Delete software client").
					WithHttpClient(secureClient).
					ClientDelete(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					Build(),
			).
			Build(),
		NewBuilder(
			"DCR-006",
			"I should not be able to retrieve a registered software if I send invalid credentials",
			specLinkRetrieveSoftware,
		).
			TestCase(
				NewTestCaseBuilder("Register software client").
					WithHttpClient(secureClient).
					GenerateSignedClaims(authoriserBuilder).
					PostClientRegister(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeCreated().
					ParseClientRegisterResponse(authoriserBuilder).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve software client with invalid credentials should not succeed").
					WithHttpClient(secureClient).
					SetInvalidGrantToken().
					ClientRetrieve(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeUnauthorized().
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve client credentials grant").
					WithHttpClient(secureClient).
					GetClientCredentialsGrant(cfg.OpenIDConfig.TokenEndpoint).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Delete software client").
					WithHttpClient(secureClient).
					ClientDelete(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					Build(),
			).
			Build(),
	}

	return NewManifest("DCR32", "1.0", scenarios)
}
