package compliant

import (
	"fmt"
	"net/http"
	"time"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/step"
	"github.com/dgrijalva/jwt-go"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/schema"
)

// nolint:lll
const (
	specLinkDiscovery        = "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2#DynamicClientRegistration-v3.2-Discovery"
	specLinkRegisterSoftware = "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2#DynamicClientRegistration-v3.2-POST/register"
	specLinkDeleteSoftware   = "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2#DynamicClientRegistration-v3.2-DELETE/register/{ClientId}"
	specLinkRetrieveSoftware = "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2#DynamicClientRegistration-v3.2-GET/register/{ClientId}"
	specLinkUpdateSoftware   = "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2#DynamicClientRegistration-v3.2-PUT/register/{ClientId}"
)

func NewDCR32(cfg DCR32Config) (Manifest, error) {
	secureClient := cfg.SecureClient
	authoriserBuilder := cfg.AuthoriserBuilder
	validator := cfg.SchemaValidator

	scenarios := Scenarios{
		DCR32ValidateOIDCConfigRegistrationURL(cfg),
		DCR32CreateSoftwareClient(cfg, secureClient, authoriserBuilder),
		DCR32DeleteSoftwareClient(cfg, secureClient, authoriserBuilder),
		DCR32CreateInvalidRegistrationRequest(cfg, secureClient, authoriserBuilder),
		DCR32RetrieveSoftwareClient(cfg, secureClient, authoriserBuilder, validator),
		DCR32RetrieveWithInvalidCredentials(cfg, secureClient, authoriserBuilder),
		DCR32UpdateSoftwareClient(cfg, secureClient, authoriserBuilder),
		DCR32UpdateSoftwareClientWithWrongId(cfg, secureClient, authoriserBuilder),
		DCR32RetrieveSoftwareClientWrongId(cfg, secureClient, authoriserBuilder),
		DCR32RegisterSoftwareWrongResponseType(cfg, secureClient, authoriserBuilder),
	}

	return NewManifest("DCR32", "1.0", scenarios)
}

func DCR32ValidateOIDCConfigRegistrationURL(cfg DCR32Config) Scenario {
	return NewBuilder(
		"DCR-001",
		"Validate OIDC Config Registration URL",
		specLinkDiscovery,
	).TestCase(
		NewTestCaseBuilder("Validate Registration URL").
			ValidateRegistrationEndpoint(cfg.OpenIDConfig.RegistrationEndpoint).
			Build(),
	).Build()
}

func DCR32CreateSoftwareClient(
	cfg DCR32Config,
	secureClient *http.Client,
	authoriserBuilder auth.AuthoriserBuilder,
) Scenario {
	return NewBuilder(
		"DCR-002",
		"Dynamically create a new software client",
		specLinkRegisterSoftware,
	).
		TestCase(DCR32CreateSoftwareClientTestCases(cfg, secureClient, authoriserBuilder)...).
		TestCase(DCR32DeleteSoftwareClientTestCase(cfg, secureClient)).
		Build()
}

func DCR32CreateSoftwareClientTestCases(
	cfg DCR32Config,
	secureClient *http.Client,
	authoriserBuilder auth.AuthoriserBuilder,
) []TestCase {
	return []TestCase{
		NewTestCaseBuilder("Register software client").
			WithHttpClient(secureClient).
			GenerateSignedClaims(authoriserBuilder).
			PostClientRegister(cfg.OpenIDConfig.RegistrationEndpointAsString()).
			AssertStatusCodeCreated().
			ParseClientRegisterResponse(authoriserBuilder).
			Build(),
		NewTestCaseBuilder("Retrieve client credentials grant").
			WithHttpClient(secureClient).
			GetClientCredentialsGrant(cfg.OpenIDConfig.TokenEndpoint).
			Build(),
	}
}

func DCR32DeleteSoftwareClientTestCase(
	cfg DCR32Config,
	secureClient *http.Client,
) TestCase {
	name := "Delete software client"
	if !cfg.DeleteImplemented {
		return NewTestCase(
			fmt.Sprintf("(SKIP Delete endpoint not implemented) %s", name),
			[]step.Step{},
		)
	}
	return NewTestCaseBuilder(name).
		WithHttpClient(secureClient).
		ClientDelete(cfg.OpenIDConfig.RegistrationEndpointAsString()).
		Build()
}

func DCR32DeleteSoftwareClient(
	cfg DCR32Config,
	secureClient *http.Client,
	authoriserBuilder auth.AuthoriserBuilder,
) Scenario {
	id := "DCR-003"
	name := "Delete software is supported"

	if !cfg.DeleteImplemented {
		return NewBuilder(
			id,
			fmt.Sprintf("(SKIP Delete endpoint not implemented) %s", name),
			specLinkDeleteSoftware,
		).Build()
	}

	return NewBuilder(
		id,
		name,
		specLinkDeleteSoftware,
	).
		TestCase(DCR32CreateSoftwareClientTestCases(cfg, secureClient, authoriserBuilder)...).
		TestCase(DCR32DeleteSoftwareClientTestCase(cfg, secureClient)).
		TestCase(
			NewTestCaseBuilder("Retrieve delete software client should fail").
				WithHttpClient(secureClient).
				ClientRetrieve(cfg.OpenIDConfig.RegistrationEndpointAsString()).
				AssertStatusCodeUnauthorized().
				Build(),
		).Build()
}

func DCR32CreateInvalidRegistrationRequest(
	cfg DCR32Config,
	secureClient *http.Client,
	authoriserBuilder auth.AuthoriserBuilder,
) Scenario {
	return NewBuilder(
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
						WithIssuer("foo.is/invalid"),
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
						WithIssuer(""),
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
						WithIssuer("123456789012345678901234567890"),
				).
				PostClientRegister(cfg.OpenIDConfig.RegistrationEndpointAsString()).
				AssertStatusCodeBadRequest().
				Build(),
		).
		TestCase(
			NewTestCaseBuilder("Register software client will fail with token endpoint auth method RS256").
				WithHttpClient(secureClient).
				GenerateSignedClaims(authoriserBuilder.WithTokenEndpointAuthMethod(jwt.SigningMethodRS256)).
				PostClientRegister(cfg.OpenIDConfig.RegistrationEndpointAsString()).
				AssertStatusCodeBadRequest().
				Build(),
		).Build()
}

func DCR32RetrieveSoftwareClient(
	cfg DCR32Config,
	secureClient *http.Client,
	authoriserBuilder auth.AuthoriserBuilder,
	validator schema.Validator,
) Scenario {
	return NewBuilder(
		"DCR-005",
		"Dynamically retrieve a new software client",
		specLinkRetrieveSoftware,
	).
		TestCase(DCR32CreateSoftwareClientTestCases(cfg, secureClient, authoriserBuilder)...).
		TestCase(DCR32RetrieveSoftwareClientTestCase(cfg, secureClient, validator)).
		TestCase(DCR32DeleteSoftwareClientTestCase(cfg, secureClient)).
		Build()
}

func DCR32RetrieveSoftwareClientTestCase(
	cfg DCR32Config,
	secureClient *http.Client,
	validator schema.Validator,
) TestCase {
	name := "Retrieve software client"
	if !cfg.GetImplemented {
		return NewTestCase(
			fmt.Sprintf("(SKIP Get endpoint not implemented) %s", name),
			[]step.Step{},
		)
	}
	return NewTestCaseBuilder("Retrieve software client").
		WithHttpClient(secureClient).
		ClientRetrieve(cfg.OpenIDConfig.RegistrationEndpointAsString()).
		AssertStatusCodeOk().
		AssertValidSchemaResponse(validator).
		ParseClientRetrieveResponse(cfg.OpenIDConfig.TokenEndpoint).
		Build()
}

func DCR32RetrieveWithInvalidCredentials(
	cfg DCR32Config,
	secureClient *http.Client,
	authoriserBuilder auth.AuthoriserBuilder,
) Scenario {
	return NewBuilder(
		"DCR-007",
		"I should not be able to retrieve a software client with invalid credentials",
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
			NewTestCaseBuilder("Retrieve software client with invalid credentials grant").
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
		TestCase(DCR32DeleteSoftwareClientTestCase(cfg, secureClient)).
		Build()
}

func DCR32UpdateSoftwareClient(
	cfg DCR32Config,
	secureClient *http.Client,
	authoriserBuilder auth.AuthoriserBuilder,
) Scenario {
	id := "DCR-008"
	const name = "I should be able update a registered software"

	if !cfg.PutImplemented {
		return NewBuilder(
			id,
			fmt.Sprintf("(SKIP PUT endpoint not implemented) %s", name),
			specLinkRetrieveSoftware,
		).Build()
	}

	return NewBuilder(
		id,
		name,
		specLinkUpdateSoftware,
	).
		TestCase(DCR32CreateSoftwareClientTestCases(cfg, secureClient, authoriserBuilder)...).
		TestCase(
			NewTestCaseBuilder("Update an existing software client").
				WithHttpClient(secureClient).
				GenerateSignedClaims(authoriserBuilder).
				ClientUpdate(cfg.OpenIDConfig.RegistrationEndpointAsString()).
				AssertStatusCodeOk().
				Build(),
		).
		TestCase(DCR32DeleteSoftwareClientTestCase(cfg, secureClient)).
		Build()
}

func DCR32UpdateSoftwareClientWithWrongId(
	cfg DCR32Config,
	secureClient *http.Client,
	authoriserBuilder auth.AuthoriserBuilder,
) Scenario {
	id := "DCR-009"
	const name = "When I try to update a non existing software client I should be unauthorized"

	if !cfg.PutImplemented {
		return NewBuilder(
			id,
			fmt.Sprintf("(SKIP PUT endpoint not implemented) %s", name),
			specLinkRetrieveSoftware,
		).Build()
	}

	return NewBuilder(
		id,
		name,
		specLinkUpdateSoftware,
	).
		TestCase(DCR32CreateSoftwareClientTestCases(cfg, secureClient, authoriserBuilder)...).
		TestCase(DCR32DeleteSoftwareClientTestCase(cfg, secureClient)).
		TestCase(
			NewTestCaseBuilder("Update a deleted software client").
				WithHttpClient(secureClient).
				GenerateSignedClaims(authoriserBuilder).
				ClientUpdate(cfg.OpenIDConfig.RegistrationEndpointAsString()).
				AssertStatusCodeUnauthorized().
				Build(),
		).Build()
}

func DCR32RetrieveSoftwareClientWrongId(
	cfg DCR32Config,
	secureClient *http.Client,
	authoriserBuilder auth.AuthoriserBuilder,
) Scenario {
	id := "DCR-010"
	const name = "When I try to retrieve a non existing software client I should be unauthorized"

	return NewBuilder(
		id,
		name,
		specLinkUpdateSoftware,
	).
		TestCase(DCR32CreateSoftwareClientTestCases(cfg, secureClient, authoriserBuilder)...).
		TestCase(DCR32DeleteSoftwareClientTestCase(cfg, secureClient)).
		TestCase(
			NewTestCaseBuilder("Retrieve a deleted software client").
				WithHttpClient(secureClient).
				ClientRetrieve(cfg.OpenIDConfig.RegistrationEndpointAsString()).
				AssertStatusCodeUnauthorized().
				Build(),
		).Build()
}

func DCR32RegisterSoftwareWrongResponseType(
	cfg DCR32Config,
	secureClient *http.Client,
	authoriserBuilder auth.AuthoriserBuilder,
) Scenario {
	id := "DCR-011"
	const name = "When I try to register a software with invalid response_types it should be fail"

	return NewBuilder(
		id,
		name,
		specLinkRegisterSoftware,
	).
		TestCase(
			NewTestCaseBuilder("Register software client").
				WithHttpClient(secureClient).
				GenerateSignedClaims(
					authoriserBuilder.WithResponseTypes([]string{"id_token", "token"}),
				).
				PostClientRegister(cfg.OpenIDConfig.RegistrationEndpointAsString()).
				AssertStatusCodeBadRequest().
				ParseClientRegisterResponse(authoriserBuilder).
				Build(),
		).
		Build()
}
