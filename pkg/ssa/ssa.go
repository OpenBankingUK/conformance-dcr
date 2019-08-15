package ssa

import (
	"github.com/dgrijalva/jwt-go"
)

// SSA is a software statement assertion
// It is an implementation of an [RFC7591] software statement, signed by the OpenBanking Directory.
type SSA struct {
	// RFC7591 header fields
	Typ string `json:"typ"` // MUST be set to JWT
	Alg string `json:"alg"` // MUST be set to ES256 or PS256

	// The kid will be kept the same as the "x5t" parameter.
	// (X.509 Certificate SHA-1 Thumbprint) of the signing certificate.
	Kid string `json:"kid"`

	// RFC7591 payload fields
	jwt.StandardClaims

	// OB SSA fields - Software Metadata
	// Unique ID for TPP client software
	SoftwareID string `json:"software_id"`

	// Requested additional field to avoid certificate check
	SoftwareEnvironment string `json:"software_environment"`

	// ASPSP Requested additional field to indicate that this software
	// is "Test" or "Live" the default is "Live". Impact and support
	// for "Test" software is up to the ASPSP.
	SoftwareMode string `json:"software_mode"`

	// The Client ID registered at OB used to access OB resources
	SoftwareClientID string `json:"software_client_id"`

	// Human-readable Software Name
	SoftwareClientName string `json:"software_client_name"`

	// Human-readable detailed description of the client
	SoftwareClientDescription string `json:"software_client_description"`

	// The version number of the software should a TPP choose to register and / or maintain it
	SoftwareVersion string `json:"software_version"`

	// The website or resource root uri
	SoftwareClientURI string `json:"software_client_uri"`

	// Contains all active signing and network certs for the software
	SoftwareJWKSEndpoint string `json:"software_jwks_endpoint"`

	// Contains all revoked signing and network certs for the software
	SoftwareJWKSRevokedEndpoint string `json:"software_jwks_revoked_endpoint"`

	// Link to the TPP logo. Note, ASPSPs are not obliged
	// to display images hosted by third parties
	SoftwareLogoURI string `json:"software_logo_uri"`

	// A potential reference to a fourth party if the TPP is registering
	// a software statement or acting on behalf of another
	SoftwareOnBehalfOfOrg string `json:"software_on_behalf_of_org"`

	// A link to the software's policy page
	SoftwarePolicyURI string `json:"sofware_policy_uri"`

	// Registered client callback endpoints as registered with Open Banking
	SoftwareRedirectURIs []string `json:"software_redirect_uris"`

	// A multi value list of PSD2 roles that this software is authorized to perform.
	SoftwareRoles []string `json:"software_roles"`

	// A link to the software's terms of service page
	SoftwareTermsOfServiceURI string `json:"software_tos_uri"`

	// OB SSA fields - Organisation Metadata
	// Authorisations granted to the organsiation by an NCA
	OrganisationCompetentAuthorityClaims string `json:"organisation_competent_authority_claims"`

	// Included to cater for voluntary withdrawal from OB scenarios
	OrganisationStatus string `json:"org_status"`

	// The Unique TPP or ASPSP ID held by OpenBanking.
	OrganisationID string `json:"org_id"`

	// Legal Entity Identifier or other known organisation name
	OrganisationName string `json:"org_name"`

	// JSON array of objects containing a triplet of name, email, and phone number
	OrganisationContacts []string `json:"org_contacts"`

	// Contains all active signing and network certs for the organisation
	OrganisationJWKSEndpoint string `json:"org_jwks_endpoint"`

	// Contains all revoked signing and network certs for the organisation
	OrganisationJWKSRevokedEndpoint string `json:"org_jwks_revoked_endpoint"`

	// A link to the OB registries terms of service page
	OBRegistryTermsOfService string `json:"ob_registry_tos"`
}
