package schema

import (
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"io"
	"regexp"
)

type responseValidator32 struct{}

func (v responseValidator32) Validate(data io.Reader) []Failure {
	var failures []Failure

	var registrationResponse OBClientRegistrationResponseSchema32
	if err := json.NewDecoder(data).Decode(&registrationResponse); err != nil {
		return []Failure{Failure(err.Error())}
	}

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
			validation.Each(validation.In("client_credentials", "authorization_code", "refresh_token")),
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
			validation.In("web", "mobile"),
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
			validation.Length(1, 128),
		),
	}
	failures = append(failures, toFailures(errs)...)

	if registrationResponse.TokenEndpointAuthMethod != nil {
		failures = append(failures, validateTokenEndpointMethod32(registrationResponse)...)
	}

	return failures
}

func validateTokenEndpointMethod32(r OBClientRegistrationResponseSchema32) []Failure {
	var failures []Failure

	if equalsAny(r.TokenEndpointAuthMethod, "private_key_jwt", "client_secret_jwt") {
		if emptyTokenEndpointAuthSignAlg(r) {
			msg := "token_endpoint_auth_signing_alg MUST be specified if token_endpoint_auth_method " +
				"is private_key_jwt or client_secret_jwt"
			failures = append(failures, Failure(msg))
		}
	}

	if equals(r.TokenEndpointAuthMethod, "tls_client_auth") {
		if emptyTLSClientAuthSubjectDn(r) {
			msg := "tls_client_auth_subject_dn MUST be set if token_endpoint_auth_method is set to tls_client_auth"
			failures = append(failures, Failure(msg))
		}
	}

	return failures
}

func equals(value *string, compare string) bool {
	return equalsAny(value, compare)
}

func equalsAny(value *string, compare ...string) bool {
	if value == nil {
		return false
	}
	for _, item := range compare {
		if item == *value {
			return true
		}
	}
	return false
}

func emptyTokenEndpointAuthSignAlg(r OBClientRegistrationResponseSchema32) bool {
	return r.TokenEndpointAuthSignAlg == nil || *r.TokenEndpointAuthSignAlg == ""
}

func emptyTLSClientAuthSubjectDn(r OBClientRegistrationResponseSchema32) bool {
	return r.TLSClientAuthSubjectDn == nil || *r.TLSClientAuthSubjectDn == ""
}

// https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2
type OBClientRegistrationResponseSchema32 struct {
	ClientID                 *string  `json:"client_id"`
	ClientSecret             *string  `json:"client_secret"`
	ClientIdIssuedAt         *int     `json:"client_id_issued_at"`
	ClientSecretExpiresAt    *int     `json:"client_secret_expires_at"`
	RedirectURIs             []string `json:"redirect_uris"`
	TokenEndpointAuthMethod  *string  `json:"token_endpoint_auth_method"`
	GrantTypes               []string `json:"grant_types"`
	ResponseTypes            []string `json:"response_types"`
	SoftwareId               *string  `json:"software_id"`
	Scope                    *string  `json:"scope"`
	SoftwareStatement        *string  `json:"software_statement"`
	ApplicationType          *string  `json:"application_type"`
	IdTokenSignedResponseAlg *string  `json:"id_token_signed_response_alg"`
	RequestObjSigningAlg     *string  `json:"request_object_signing_alg"`
	TokenEndpointAuthSignAlg *string  `json:"token_endpoint_auth_signing_alg"`
	TLSClientAuthSubjectDn   *string  `json:"tls_client_auth_subject_dn"`
}

func toFailures(errs map[string]error) []Failure {
	var failures []Failure
	for key, err := range errs {
		if err != nil {
			msg := fmt.Sprintf("%s: %s", key, err.Error())
			failures = append(failures, Failure(msg))
		}
	}
	return failures
}
