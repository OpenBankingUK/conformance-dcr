package compliant

import (
	"net/http"
	"testing"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/schema"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/step"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDCR32(t *testing.T) {
	manifest, err := NewDCR32(DCR32Config{SSA: "ssa"})
	require.NoError(t, err)

	assert.Equal(t, "1.0", manifest.Version())
	assert.Equal(t, "DCR32", manifest.Name())
	assert.Equal(t, 10, len(manifest.Scenarios()))
}

func TestDCR32ValidateOIDCConfigRegistrationURL(t *testing.T) {
	scenario := DCR32ValidateOIDCConfigRegistrationURL(
		DCR32Config{},
	)

	assert.Equal(t, "DCR-001", scenario.Id())
	name := "Validate OIDC Config Registration URL"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkDiscovery, scenario.Spec())
}

func TestDCR32CreateSoftwareClient(t *testing.T) {
	scenario := DCR32CreateSoftwareClient(
		DCR32Config{DeleteImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
		&[]string{},
	)

	assert.Equal(t, "DCR-002", scenario.Id())
	name := "Dynamically create a new software client"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkRegisterSoftware, scenario.Spec())
}

func TestDCR32DeleteSoftwareClient(t *testing.T) {
	scenario := DCR32DeleteSoftwareClient(
		DCR32Config{DeleteImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
		&[]string{},
	)

	assert.Equal(t, "DCR-003", scenario.Id())
	name := "Delete software is supported"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkDeleteSoftware, scenario.Spec())
}

func TestDCR32DeleteSoftwareClient_DeleteNotImplemented(t *testing.T) {
	scenario := DCR32DeleteSoftwareClient(
		DCR32Config{DeleteImplemented: false},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
		&[]string{},
	)

	result := scenario.Run()

	assert.Equal(t, "DCR-003", scenario.Id())
	name := "(SKIP Delete endpoint not implemented) Delete software is supported"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkDeleteSoftware, scenario.Spec())
	assert.False(t, result.Fail())
}

func TestDCR32CreateInvalidRegistrationRequest(t *testing.T) {
	scenario := DCR32CreateInvalidRegistrationRequest(
		DCR32Config{GetImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
		&[]string{},
	)

	assert.Equal(t, "DCR-004", scenario.Id())
	name := "Dynamically create a new software client will fail on invalid registration request"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkRegisterSoftware, scenario.Spec())
}

func TestDCR32RetrieveSoftwareClient(t *testing.T) {
	validator, err := schema.NewValidator("3.2")
	require.NoError(t, err)
	scenario := DCR32RetrieveSoftwareClient(
		DCR32Config{GetImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
		validator,
		&[]string{},
	)

	assert.Equal(t, "DCR-005", scenario.Id())
	assert.Equal(t, "Dynamically retrieve a new software client", scenario.Name())
	assert.Equal(t, specLinkRetrieveSoftware, scenario.Spec())
}

func TestTCRetrieveSoftwareClient(t *testing.T) {
	validator, err := schema.NewValidator("3.2")
	require.NoError(t, err)
	tc := DCR32RetrieveSoftwareClientTestCase(
		DCR32Config{GetImplemented: true},
		&http.Client{},
		validator,
	)

	result := tc.Run(step.NewContext())

	assert.Equal(t, "Retrieve software client", result.Name)
	assert.True(t, result.Fail())
}

func TestTCRetrieveSoftwareClient_GetNotImplemented(t *testing.T) {
	validator, err := schema.NewValidator("3.2")
	require.NoError(t, err)
	tc := DCR32RetrieveSoftwareClientTestCase(
		DCR32Config{GetImplemented: false},
		&http.Client{},
		validator,
	)

	result := tc.Run(step.NewContext())

	assert.Equal(t, "(SKIP Get endpoint not implemented) Retrieve software client", result.Name)
	assert.Equal(t, step.Results(nil), result.Results)
	assert.False(t, result.Fail())
}

func TestDCR32RetrieveWithInvalidCredentials(t *testing.T) {
	scenario := DCR32RetrieveWithInvalidCredentials(
		DCR32Config{},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
		&[]string{},
	)

	assert.Equal(t, "DCR-007", scenario.Id())
	name := "I should not be able to retrieve a software client with invalid credentials"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkRetrieveSoftware, scenario.Spec())
}

func TestDCR32UpdateSoftwareClient(t *testing.T) {
	scenario := DCR32UpdateSoftwareClient(
		DCR32Config{PutImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
		&[]string{},
	)

	assert.Equal(t, "DCR-008", scenario.Id())
	name := "I should be able update a registered software"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkUpdateSoftware, scenario.Spec())
}

func TestDCR32UpdateSoftwareClientDisabled(t *testing.T) {
	scenario := DCR32UpdateSoftwareClient(
		DCR32Config{PutImplemented: false},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
		&[]string{},
	)

	assert.Equal(t, "DCR-008", scenario.Id())
	name := "(SKIP PUT endpoint not implemented) I should be able update a registered software"
	assert.Equal(t, name, scenario.Name())
}

func TestDCR32UpdateWrongId(t *testing.T) {
	scenario := DCR32UpdateSoftwareClientWithWrongId(
		DCR32Config{PutImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
		&[]string{},
	)

	assert.Equal(t, "DCR-009", scenario.Id())
	name := "When I try to update a non existing software client I should be unauthorized"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkUpdateSoftware, scenario.Spec())
}

func TestDCR32RetrieveWrongId(t *testing.T) {
	scenario := DCR32RetrieveSoftwareClientWrongId(
		DCR32Config{},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
		&[]string{},
	)

	assert.Equal(t, "DCR-010", scenario.Id())
	name := "When I try to retrieve a non existing software client I should be unauthorized"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkUpdateSoftware, scenario.Spec())
}

func TestDCR32RegisterWrongResponseTypes(t *testing.T) {
	scenario := DCR32RegisterSoftwareWrongResponseType(
		DCR32Config{},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
		&[]string{},
	)

	assert.Equal(t, "DCR-011", scenario.Id())
	name := "When I try to register a software with invalid response_types it should be fail"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkRegisterSoftware, scenario.Spec())
}
