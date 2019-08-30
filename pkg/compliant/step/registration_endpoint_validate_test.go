package step

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidateRegistrationEndpoint_ReturnsSuccessfullResult(t *testing.T) {
	url := "http://x.org/api/register"
	registrationEndpointStep := NewValidateRegistrationEndpoint(&url)
	result := registrationEndpointStep.Run(NewContext())
	assert.True(t, result.Pass)
	assert.Equal(t, result.Name, "Registration Endpoint Validate")
}

func TestNewValidateRegistrationEndpoint_ReturnsFailureResultOnInvalidEndpoint(t *testing.T) {
	url := "foo/bar"
	registrationEndpointStep := NewValidateRegistrationEndpoint(&url)
	result := registrationEndpointStep.Run(NewContext())
	assert.False(t, result.Pass)
	assert.Equal(t, result.Name, "Registration Endpoint Validate")
	assert.Contains(t, result.FailReason, "registration endpoint foo/bar is invalid")
}

func TestNewValidateRegistrationEndpoint_ReturnsFailureResultOnBlankEndpoint(t *testing.T) {
	registrationEndpointStep := NewValidateRegistrationEndpoint(nil)
	result := registrationEndpointStep.Run(NewContext())
	assert.False(t, result.Pass)
	assert.Equal(t, result.Name, "Registration Endpoint Validate")
	assert.Contains(t, result.FailReason, "registration endpoint is missing")
}
