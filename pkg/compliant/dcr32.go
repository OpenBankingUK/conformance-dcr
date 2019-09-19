package compliant

import (
	"net/http"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
)

func NewDCR32(
	cfg DCR32Config,
	secureClient *http.Client,
	authoriserBuilder auth.AuthoriserBuilder,
) Scenarios {
	return Scenarios{
		NewBuilder("Validate OIDC Config Registration URL").
			TestCase(
				NewTestCaseBuilder("Validate Registration URL").
					ValidateRegistrationEndpoint(cfg.OpenIDConfig.RegistrationEndpoint).
					Build(),
			).
			Build(),
		NewBuilder("Dynamically create a new software client").
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
		NewBuilder("Dynamically create a new software client will fail on invalid registration request").
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
								openid.Configuration{
									RegistrationEndpoint:              cfg.OpenIDConfig.RegistrationEndpoint,
									TokenEndpoint:                     cfg.OpenIDConfig.TokenEndpoint,
									Issuer:                            "foo.is/invalid",
									ObjectSignAlgSupported:            cfg.OpenIDConfig.ObjectSignAlgSupported,
									TokenEndpointAuthMethodsSupported: cfg.OpenIDConfig.TokenEndpointAuthMethodsSupported,
								},
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
								openid.Configuration{
									RegistrationEndpoint:              cfg.OpenIDConfig.RegistrationEndpoint,
									TokenEndpoint:                     cfg.OpenIDConfig.TokenEndpoint,
									Issuer:                            "",
									ObjectSignAlgSupported:            cfg.OpenIDConfig.ObjectSignAlgSupported,
									TokenEndpointAuthMethodsSupported: cfg.OpenIDConfig.TokenEndpointAuthMethodsSupported,
								},
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
								openid.Configuration{
									RegistrationEndpoint:              cfg.OpenIDConfig.RegistrationEndpoint,
									TokenEndpoint:                     cfg.OpenIDConfig.TokenEndpoint,
									Issuer:                            "123456789012345678901234567890",
									ObjectSignAlgSupported:            cfg.OpenIDConfig.ObjectSignAlgSupported,
									TokenEndpointAuthMethodsSupported: cfg.OpenIDConfig.TokenEndpointAuthMethodsSupported,
								},
							),
					).
					PostClientRegister(cfg.OpenIDConfig.RegistrationEndpointAsString()).
					AssertStatusCodeBadRequest().
					Build(),
			).
			Build(),
		NewBuilder("Dynamically retrieve a new software client").
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
		NewBuilder("I should not be able to retrieve a registered software if I send invalid credentials").
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
}
