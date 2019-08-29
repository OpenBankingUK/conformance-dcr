package step

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidateRegistrationEndpoint_Run_ReturnsSuccessfullResult(t *testing.T) {
	registrationEndpointStep := NewValidateRegistrationEndpoint("http://x.org/api/register")
	result := registrationEndpointStep.Run(NewContext())
	assert.True(t, result.Pass)
	assert.Equal(t, result.Name, "Registration Endpoint Validate")
}

func TestNewValidateRegistrationEndpoint_Run_ReturnsFailureResultOnInvalidEndpoint(t *testing.T) {
	registrationEndpointStep := NewValidateRegistrationEndpoint("foo/bar")
	result := registrationEndpointStep.Run(NewContext())
	assert.False(t, result.Pass)
	assert.Equal(t, result.Name, "Registration Endpoint Validate")
}

func TestNewValidateRegistrationEndpoint_Run_ReturnsFailureResultOnBlankEndpoint(t *testing.T) {
	registrationEndpointStep := NewValidateRegistrationEndpoint("")
	result := registrationEndpointStep.Run(NewContext())
	assert.False(t, result.Pass)
	assert.Equal(t, result.Name, "Registration Endpoint Validate")
}
