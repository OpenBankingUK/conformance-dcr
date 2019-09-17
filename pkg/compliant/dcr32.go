package compliant

import (
	"net/http"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/config"
)

func NewDCR32(
	openIDConfig openid.Configuration,
	secureClient *http.Client,
	authoriserBuilder auth.AuthoriserBuilder,
	cfg config.Config,
) Scenarios {
	return Scenarios{
		NewBuilder("Validate OIDC Config Registration URL").
			TestCase(
				NewTestCaseBuilder("Validate Registration URL").
					ValidateRegistrationEndpoint(openIDConfig.RegistrationEndpoint).
					Build(),
			).
			Build(),
		NewBuilder("Dynamically create a new software client").
			TestCase(
				NewTestCaseBuilder("Register software client").
					WithHttpClient(secureClient).
					GenerateSignedClaims(
						authoriserBuilder.
							WithOpenIDConfig(openIDConfig).
							WithSSA(cfg.SSA).
							WithKID(cfg.Kid).
							WithClientID(cfg.ClientId).
							WithRedirectURIs(cfg.RedirectURIs).
							WithPrivateKey(cfg.PrivateKeyBytes).
							WithJwtExpiration(time.Hour).
							Build(),
					).
					PostClientRegister(openIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeCreated().
					ParseClientRegisterResponse(
						authoriserBuilder.
							WithOpenIDConfig(openIDConfig).
							WithSSA(cfg.SSA).
							WithKID(cfg.Kid).
							WithClientID(cfg.ClientId).
							WithRedirectURIs(cfg.RedirectURIs).
							WithPrivateKey(cfg.PrivateKeyBytes).
							WithJwtExpiration(time.Hour).
							Build(),
					).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve client credentials grant").
					WithHttpClient(secureClient).
					GetClientCredentialsGrant(openIDConfig.TokenEndpoint).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Delete software client").
					WithHttpClient(secureClient).
					ClientDelete(openIDConfig.RegistrationEndpointAsString()).
					Build(),
			).
			Build(),
		NewBuilder("Dynamically create a new software client will fail on invalid registration request").
			TestCase(
				NewTestCaseBuilder("Register software client fails on expired claims").
					WithHttpClient(secureClient).
					GenerateSignedClaims(
						authoriserBuilder.
							WithOpenIDConfig(openIDConfig).
							WithSSA(cfg.SSA).
							WithKID(cfg.Kid).
							WithClientID(cfg.ClientId).
							WithRedirectURIs(cfg.RedirectURIs).
							WithPrivateKey(cfg.PrivateKeyBytes).
							WithJwtExpiration(-time.Hour).
							Build(),
					).
					PostClientRegister(openIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeBadRequest().
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Register software client fails on invalid issuer").
					WithHttpClient(secureClient).
					GenerateSignedClaims(
						authoriserBuilder.
							WithOpenIDConfig(
								openid.Configuration{
									RegistrationEndpoint:              openIDConfig.RegistrationEndpoint,
									TokenEndpoint:                     openIDConfig.TokenEndpoint,
									Issuer:                            "foo.is/invalid",
									ObjectSignAlgSupported:            openIDConfig.ObjectSignAlgSupported,
									TokenEndpointAuthMethodsSupported: openIDConfig.TokenEndpointAuthMethodsSupported,
								},
							).
							WithSSA(cfg.SSA).
							WithKID(cfg.Kid).
							WithClientID(cfg.ClientId).
							WithRedirectURIs(cfg.RedirectURIs).
							WithPrivateKey(cfg.PrivateKeyBytes).
							WithJwtExpiration(-time.Hour).
							Build(),
					).
					PostClientRegister(openIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeBadRequest().
					Build(),
			).
			Build(),
		NewBuilder("Dynamically retrieve a new software client").
			TestCase(
				NewTestCaseBuilder("Register software client").
					WithHttpClient(secureClient).
					GenerateSignedClaims(
						authoriserBuilder.
							WithOpenIDConfig(openIDConfig).
							WithSSA(cfg.SSA).
							WithKID(cfg.Kid).
							WithClientID(cfg.ClientId).
							WithRedirectURIs(cfg.RedirectURIs).
							WithPrivateKey(cfg.PrivateKeyBytes).
							WithJwtExpiration(time.Hour).
							Build(),
					).
					PostClientRegister(openIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeCreated().
					ParseClientRegisterResponse(
						authoriserBuilder.
							WithOpenIDConfig(openIDConfig).
							WithSSA(cfg.SSA).
							WithKID(cfg.Kid).
							WithClientID(cfg.ClientId).
							WithRedirectURIs(cfg.RedirectURIs).
							WithPrivateKey(cfg.PrivateKeyBytes).
							WithJwtExpiration(time.Hour).
							Build(),
					).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve client credentials grant").
					WithHttpClient(secureClient).
					GetClientCredentialsGrant(openIDConfig.TokenEndpoint).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve software client").
					WithHttpClient(secureClient).
					ClientRetrieve(openIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeOk().
					ParseClientRetrieveResponse(openIDConfig.TokenEndpoint).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Delete software client").
					WithHttpClient(secureClient).
					ClientDelete(openIDConfig.RegistrationEndpointAsString()).
					Build(),
			).
			Build(),
		NewBuilder("I should not be able to retrieve a registered software if I send invalid credentials").
			TestCase(
				NewTestCaseBuilder("Register software client").
					WithHttpClient(secureClient).
					GenerateSignedClaims(
						authoriserBuilder.
							WithOpenIDConfig(openIDConfig).
							WithSSA(cfg.SSA).
							WithKID(cfg.Kid).
							WithClientID(cfg.ClientId).
							WithRedirectURIs(cfg.RedirectURIs).
							WithPrivateKey(cfg.PrivateKeyBytes).
							WithJwtExpiration(time.Hour).
							Build(),
					).
					PostClientRegister(openIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeCreated().
					ParseClientRegisterResponse(
						authoriserBuilder.
							WithOpenIDConfig(openIDConfig).
							WithSSA(cfg.SSA).
							WithKID(cfg.Kid).
							WithClientID(cfg.ClientId).
							WithRedirectURIs(cfg.RedirectURIs).
							WithPrivateKey(cfg.PrivateKeyBytes).
							WithJwtExpiration(time.Hour).
							Build(),
					).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve software client with invalid credentials should not succeed").
					WithHttpClient(secureClient).
					SetInvalidGrantToken().
					ClientRetrieve(openIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeUnauthorized().
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Retrieve client credentials grant").
					WithHttpClient(secureClient).
					GetClientCredentialsGrant(openIDConfig.TokenEndpoint).
					Build(),
			).
			TestCase(
				NewTestCaseBuilder("Delete software client").
					WithHttpClient(secureClient).
					ClientDelete(openIDConfig.RegistrationEndpointAsString()).
					Build(),
			).
			Build(),
	}
}
