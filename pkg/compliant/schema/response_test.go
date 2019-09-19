package schema

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewValidator_Return32(t *testing.T) {
	validator, err := NewValidator("3.2")

	require.NoError(t, err)
	assert.IsType(t, responseValidator32{}, validator)
}

func TestNewValidator_ReturnErrorForNotSupportedSpecVersion(t *testing.T) {
	validator, err := NewValidator("3.1")

	assert.EqualError(t, err, "unknown spec version to validate schema 3.1")
	assert.Nil(t, validator)
}
