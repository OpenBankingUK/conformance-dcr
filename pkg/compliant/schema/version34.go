package schema

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"io"
	"regexp"
)

type responseValidator34 struct{}

func (v responseValidator34) Validate(data io.Reader) []Failure {
	var failures []Failure

	var registrationResponse OBClientRegistrationResponseSchema32
	if err := json.NewDecoder(data).Decode(&registrationResponse); err != nil {
		return []Failure{Failure(err.Error())}
	}

	// The spec pattern ^(?=.{1,512}$)[Cc][Nn]=[^,]+(?:,[^,\s]+=[^,]+)*$ uses a lookahead
	// which is not supported by Go's RE2 engine. The length constraint is handled separately
	// via validation.Length(1, 512) and the structural pattern is checked here.
	tlsSubjectDnPattern := regexp.MustCompile(`^[Cc][Nn]=[^,]+(?:,[^,\s]+=[^,]+)*$`)

	errs := validation.Errors{
		"client_id": validation.Validate(
			registrationResponse.ClientID,
			validation.Required,
			validation.Length(1, 36),
		),
		"client_secret": validation.Validate(
			registrationResponse.ClientSecret,
			validation.Length(1, 36),
		),
		"redirect_uris": validation.Validate(
			registrationResponse.RedirectURIs,
			validation.Required,
			validation.Each(validation.Length(1, 256), isOBURLValidationRule()),
		),
		"token_endpoint_auth_method": validation.Validate(
			registrationResponse.TokenEndpointAuthMethod,
			validation.Required,
			validation.In(
				"private_key_jwt",
				"client_secret_jwt",
				"client_secret_basic",
				"client_secret_post",
				"tls_client_auth",
			),
		),
		"grant_types": validation.Validate(
			registrationResponse.GrantTypes,
			validation.Required,
			validation.Each(
				validation.In(
					"client_credentials",
					"authorization_code",
					"refresh_token",
					"urn:openid:params:grant-type:ciba",
					"urn:ietf:params:oauth:grant-type:jwt-bearer",
				),
			),
		),
		"response_types": validation.Validate(
			registrationResponse.ResponseTypes,
			validation.Each(validation.In("code", "code id_token")),
		),
		"software_id": validation.Validate(
			registrationResponse.SoftwareId,
			validation.Match(regexp.MustCompile("^[0-9a-zA-Z]{1,22}$")),
		),
		"scope": validation.Validate(
			registrationResponse.Scope,
			validation.Required,
			validation.Length(1, 256),
		),
		"software_statement": validation.Validate(
			registrationResponse.SoftwareStatement,
			validation.Required,
		),
		"application_type": validation.Validate(
			registrationResponse.ApplicationType,
			validation.Required,
			validation.In("web", "native"),
		),
		"id_token_signed_response_alg": validation.Validate(
			registrationResponse.IdTokenSignedResponseAlg,
			validation.Required,
			validation.Length(1, 5),
		),
		"request_object_signing_alg": validation.Validate(
			registrationResponse.RequestObjSigningAlg,
			validation.Required,
			validation.Length(1, 5),
		),
		"token_endpoint_auth_signing_alg": validation.Validate(
			registrationResponse.TokenEndpointAuthSignAlg,
			validation.Length(1, 5),
		),
		"tls_client_auth_subject_dn": validation.Validate(
			registrationResponse.TLSClientAuthSubjectDn,
			validation.Length(1, 512),
			validation.Match(tlsSubjectDnPattern),
		),
	}
	failures = append(failures, toFailures(errs)...)

	if registrationResponse.TokenEndpointAuthMethod != nil {
		failures = append(failures, validateTokenEndpointMethod32(registrationResponse)...)
	}

	return failures
}
