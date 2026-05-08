package schema

import (
	"bytes"
	"fmt"
	"strings"

	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseValidator34_ValidateInvalidPayload(t *testing.T) {
	validator := responseValidator34{}
	reader := bytes.NewReader([]byte(`{`))

	failures := validator.Validate(reader)

	assert.Len(t, failures, 1)
}

func TestResponseValidator34_ValidateEmpty(t *testing.T) {
	validator := responseValidator34{}
	reader := bytes.NewReader([]byte(`{}`))

	failures := validator.Validate(reader)

	assert.Len(t, failures, 9)
}

func TestResponseValidator34_ValidateResponse(t *testing.T) {
	validator := responseValidator34{}
	reader, err := os.Open("testdata/response34.json")
	require.NoError(t, err)

	failures := validator.Validate(reader)

	assert.Len(t, failures, 0)
}

func TestResponseValidator34_ApplicationType_Native_Valid(t *testing.T) {
	appType := "native"
	response := minimalValidResponse34()
	response.ApplicationType = &appType
	reader := bytes.NewReader(marshalResponse(t, response))

	failures := responseValidator34{}.Validate(reader)

	assertNoFailureForField(t, failures, "application_type")
}

func TestResponseValidator34_ApplicationType_Mobile_Invalid(t *testing.T) {
	appType := "mobile"
	response := minimalValidResponse34()
	response.ApplicationType = &appType
	reader := bytes.NewReader(marshalResponse(t, response))

	failures := responseValidator34{}.Validate(reader)

	assertFailureForField(t, failures, "application_type")
}

func TestResponseValidator34_ApplicationType_Web_Valid(t *testing.T) {
	appType := "web"
	response := minimalValidResponse34()
	response.ApplicationType = &appType
	reader := bytes.NewReader(marshalResponse(t, response))

	failures := responseValidator34{}.Validate(reader)

	assertNoFailureForField(t, failures, "application_type")
}

func TestResponseValidator34_TLSSubjectDN_512Chars_Valid(t *testing.T) {
	// Build a DN exactly 512 chars: CN=<value>,O=<padding>
	padding := strings.Repeat("x", 512-len("CN=a,O="))
	dn := fmt.Sprintf("CN=a,O=%s", padding)
	require.Len(t, dn, 512)

	response := minimalValidResponse34()
	response.TLSClientAuthSubjectDn = &dn
	reader := bytes.NewReader(marshalResponse(t, response))

	failures := responseValidator34{}.Validate(reader)

	assertNoFailureForField(t, failures, "tls_client_auth_subject_dn")
}

func TestResponseValidator34_TLSSubjectDN_513Chars_Invalid(t *testing.T) {
	padding := strings.Repeat("x", 513-len("CN=a,O="))
	dn := fmt.Sprintf("CN=a,O=%s", padding)
	require.Len(t, dn, 513)

	response := minimalValidResponse34()
	response.TLSClientAuthSubjectDn = &dn
	reader := bytes.NewReader(marshalResponse(t, response))

	failures := responseValidator34{}.Validate(reader)

	assertFailureForField(t, failures, "tls_client_auth_subject_dn")
}

func TestResponseValidator34_TLSSubjectDN_InvalidPattern_NoCN(t *testing.T) {
	// Does not start with CN=
	dn := "O=OpenBanking,C=GB"
	response := minimalValidResponse34()
	response.TLSClientAuthSubjectDn = &dn
	reader := bytes.NewReader(marshalResponse(t, response))

	failures := responseValidator34{}.Validate(reader)

	assertFailureForField(t, failures, "tls_client_auth_subject_dn")
}

func TestResponseValidator34_TLSSubjectDN_128Chars_Valid(t *testing.T) {
	// Ensure the v3.2/v3.3 limit of 128 is also valid under v3.4
	padding := strings.Repeat("x", 128-len("CN=a,O="))
	dn := fmt.Sprintf("CN=a,O=%s", padding)
	require.Len(t, dn, 128)

	response := minimalValidResponse34()
	response.TLSClientAuthSubjectDn = &dn
	reader := bytes.NewReader(marshalResponse(t, response))

	failures := responseValidator34{}.Validate(reader)

	assertNoFailureForField(t, failures, "tls_client_auth_subject_dn")
}

func TestResponseValidator34_GrantTypes_JWTBearer_Valid(t *testing.T) {
	response := minimalValidResponse34()
	response.GrantTypes = []string{"client_credentials", "urn:ietf:params:oauth:grant-type:jwt-bearer"}
	reader := bytes.NewReader(marshalResponse(t, response))

	failures := responseValidator34{}.Validate(reader)

	assertNoFailureForField(t, failures, "grant_types")
}

// assertFailureForField checks that at least one failure mentions the given field name.
func assertFailureForField(t *testing.T, failures []Failure, field string) {
	t.Helper()
	for _, f := range failures {
		if strings.HasPrefix(string(f), field+":") {
			return
		}
	}
	t.Errorf("expected a failure for field %q but got: %v", field, failures)
}

// assertNoFailureForField checks that no failure mentions the given field name.
func assertNoFailureForField(t *testing.T, failures []Failure, field string) {
	t.Helper()
	for _, f := range failures {
		if strings.HasPrefix(string(f), field+":") {
			t.Errorf("unexpected failure for field %q: %s", field, f)
		}
	}
}
