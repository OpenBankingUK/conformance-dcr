package ssa

// SSA is a software statement assertion
// It is an implementation of an [RFC7591] software statement, signed by the OpenBanking Directory.
// For further details refer to https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2#DynamicClientRegistration-v3.2-SoftwareStatement
type SSA struct {
	// RFC7591 header fields
	SSAHeader

	// RFC7591 payload fields
	SSAPayload

	// OB SSA fields - Software Metadata
	SSASoftwareMeta

	// OB SSA fields - Organisation Metadata
	SSAOrganisationMeta
}

// SSAHeader represents SSA Header fields as defined in RFC7591
type SSAHeader struct {
	Typ string `json:"typ"` // MUST be set to JWT
	Alg string `json:"alg"` // MUST be set to ES256 or PS256
	Kid string `json:"kid"` // The kid will be kept the same as the "x5t" parameter. (X.509 Certificate SHA-1 Thumbprint) of the signing certificate.
}

// SSAPayload represents the SSA Payload fields as defined in RFC7591
type SSAPayload struct {
	Iss        string `json:"iss"`         // SSA issuer
	Iat        int64  `json:"iat"`         // Time SSA issued
	Jti        string `json:"jti"`         // JWT ID
	SoftwareID string `json:"software_id"` // Unique ID for TPP client software
}

// SSASoftwareMeta represents the SSA Payload Software Metadata
type SSASoftwareMeta struct {
	SoftwasreEnvironment        string   `json:"software_environment"`           // Requested additional field to avoid certificate check
	SoftwareMode                string   `json:"software_mode"`                  // ASPSP Requested additional field to indicate that this software is "Test" or "Live" the default is "Live". Impact and support for "Test" software is up to the ASPSP.
	SoftwareClientID            string   `json:"software_client_id"`             // The Client ID registered at OB used to access OB resources
	SoftwareClientName          string   `json:"software_client_name"`           // Human-readable Software Name
	SoftwareClientDescription   string   `json:"software_client_description"`    // Human-readable detailed description of the client
	SoftwareVersion             string   `json:"software_version"`               // The version number of the software should a TPP choose to register and / or maintain it
	SoftwareClientURI           string   `json:"software_client_uri"`            // The website or resource root uri
	SoftwareJWKSEndpoint        string   `json:"software_jwks_endpoint"`         // Contains all active signing and network certs for the software
	SoftwareJWKSRevokedEndpoint string   `json:"software_jwks_revoked_endpoint"` // Contains all revoked signing and network certs for the software
	SoftwareLogoURI             string   `json:"software_logo_uri"`              // Link to the TPP logo. Note, ASPSPs are not obliged to display images hosted by third parties
	SoftwareOnBehalfOfOrg       string   `json:"software_on_behalf_of_org"`      // A potential reference to a fourth party if the TPP is registering a software statement or acting on behalf of another
	SoftwarePolicyURI           string   `json:"sofware_policy_uri"`             // A link to the software's policy page
	SoftwareRedirectURIs        []string `json:"software_redirect_uris"`         // Registered client callback endpoints as registered with Open Banking
	SoftwareRoles               []string `json:"software_roles"`                 // A multi value list of PSD2 roles that this software is authorized to perform.
	SoftwareTermsOfServiceURI   string   `json:"software_terms_of_service_uri"`  // A link to the software's terms of service page
}

// SSAOrganisationMeta represents the SSA Payload Organisation Metadata
type SSAOrganisationMeta struct {
	OrganisationCompetentAuthorityClaims string
	OrganisationStatus                   string
	OrganisationID                       string
	OrganisationName                     string
	OrganisationContacts                 []string
	OrganisationJWKSEndpoint             string
	OrganisationJWKSRevokedEndpoint      string
	OBRegistryTermsOfService             string
}

// ParseSSA parses a SSA jwt and returns a SSA struct
// returns error in case of failure
func ParseSSA(ssa string) (SSA, error) {
	return SSA{}, nil
}
