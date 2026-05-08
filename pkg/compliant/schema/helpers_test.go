package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// minimalValidResponse34 returns an OBClientRegistrationResponseSchema32 with all fields
// required by the v3.4 validator populated with valid values.
func minimalValidResponse34() OBClientRegistrationResponseSchema32 {
	str := func(s string) *string { return &s }

	dn := "CN=TestTPP,O=OpenBanking,C=GB"
	return OBClientRegistrationResponseSchema32{
		ClientID:                 str("b424b6df-06a1-46e7-b86c-f1554b229e69"),
		RedirectURIs:             []string{"https://example.com/callback"},
		TokenEndpointAuthMethod:  str("private_key_jwt"),
		GrantTypes:               []string{"client_credentials", "authorization_code"},
		Scope:                    str("accounts openid"),
		SoftwareStatement:        str("header.payload.sig"),
		ApplicationType:          str("native"),
		IdTokenSignedResponseAlg: str("PS256"),
		RequestObjSigningAlg:     str("PS256"),
		TokenEndpointAuthSignAlg: str("PS256"),
		TLSClientAuthSubjectDn:   &dn,
	}
}

// marshalResponse JSON-encodes r and returns the bytes.
func marshalResponse(t *testing.T, r OBClientRegistrationResponseSchema32) []byte {
	t.Helper()
	b, err := json.Marshal(r)
	require.NoError(t, err)
	return b
}
