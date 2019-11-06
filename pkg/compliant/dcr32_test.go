package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/schema"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestNewDCR32(t *testing.T) {
	validator, err := schema.NewValidator("3.2")
	require.NoError(t, err)
	manifest, err := NewDCR32(
		DCR32Config{},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
		validator,
	)
	require.NoError(t, err)

	assert.Equal(t, "1.0", manifest.Version())
	assert.Equal(t, "DCR32", manifest.Name())
	assert.Equal(t, 7, len(manifest.Scenarios()))
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
	)

	assert.Equal(t, "DCR-003", scenario.Id())
	name := "Delete software statement is supported"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkDeleteSoftware, scenario.Spec())
}

func TestDCR32DeleteSoftwareClient_DeleteNotImplemented(t *testing.T) {
	scenario := DCR32DeleteSoftwareClient(
		DCR32Config{DeleteImplemented: false},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	result := scenario.Run()

	assert.Equal(t, "DCR-003", scenario.Id())
	name := "(SKIP Delete endpoint not implemented) Delete software statement is supported"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkDeleteSoftware, scenario.Spec())
	assert.False(t, result.Fail())
}

func TestDCR32CreateInvalidRegistrationRequest(t *testing.T) {
	scenario := DCR32CreateInvalidRegistrationRequest(
		DCR32Config{GetImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
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
	)

	assert.Equal(t, "DCR-005", scenario.Id())
	assert.Equal(t, "Dynamically retrieve a new software client", scenario.Name())
	assert.Equal(t, specLinkRetrieveSoftware, scenario.Spec())
}

func TestTCRetrieveSoftwareClient(t *testing.T) {
	validator, err := schema.NewValidator("3.2")
	require.NoError(t, err)
	tc := TCRetrieveSoftwareClient(
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
	tc := TCRetrieveSoftwareClient(
		DCR32Config{GetImplemented: false},
		&http.Client{},
		validator,
	)

	result := tc.Run(step.NewContext())

	assert.Equal(t, "(SKIP Get endpoint not implemented) Retrieve software client", result.Name)
	assert.Equal(t, step.Results(nil), result.Results)
	assert.False(t, result.Fail())
}

func TestDCR32RegisterWithInvalidCredentials(t *testing.T) {
	scenario := DCR32RegisterWithInvalidCredentials(
		DCR32Config{GetImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	assert.Equal(t, "DCR-006", scenario.Id())
	name := "I should not be able to retrieve a registered software if " +
		"I send invalid credentials"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkRetrieveSoftware, scenario.Spec())
}

func TestDCR32RegisterWithInvalidCredentials_GeNotImplemented(t *testing.T) {
	scenario := DCR32RegisterWithInvalidCredentials(
		DCR32Config{GetImplemented: false},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	result := scenario.Run()

	assert.Equal(t, "DCR-006", scenario.Id())
	name := "(SKIP Get endpoint not implemented) I should not be able to " +
		"retrieve a registered software if I send invalid credentials"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkRetrieveSoftware, scenario.Spec())
	assert.False(t, result.Fail())
}

func TestDCR32RetrieveWithInvalidCredentials(t *testing.T) {
	scenario := DCR32RetrieveWithInvalidCredentials(
		DCR32Config{},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	assert.Equal(t, "DCR-007", scenario.Id())
	name := "I should not be able to retrieve a registered software if " +
		"I send invalid credentials"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkRetrieveSoftware, scenario.Spec())
}
