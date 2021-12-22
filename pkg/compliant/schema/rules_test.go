package schema

import (
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
)

func TestIsURL(t *testing.T) {
	tcs := []struct {
		name  string
		url   string
		valid bool
	}{
		{
			name:  "valid url",
			url:   "https://0.0.0.0",
			valid: true,
		},
		{
			name:  "invalid schema",
			url:   "http://0.0.0.0",
			valid: false,
		},
		{
			name:  "invalid host",
			url:   "https://localhost",
			valid: false,
		},
		{
			name:  "invalid host",
			url:   "https://www.google.localhost",
			valid: false,
		},
		{
			name:  "invalid host",
			url:   "https://127.0.0.1",
			valid: false,
		},
		{
			name:  "invalid url",
			url:   string(rune(0x7f)),
			valid: false,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.valid, isOBURl(tc.url))
		})
	}
}

func TestIsOBURL(t *testing.T) {
	err := validation.Validate("http://invalid_url", isOBURLValidationRule())

	assert.EqualError(t, err, errorMessage)
}

func TestMapToFailures(t *testing.T) {
	errs := map[string]error{
		"error1": errors.New("uh-oh"),
		"error2": errors.New("ups I did it again"),
	}

	failures := toFailures(errs)

	// to not change to comparing to a slice result
	// order of the failures is NOT GUARANTEED
	assert.Len(t, failures, 2)
	assert.True(t, inSlice(failures, "error1: uh-oh"))
	assert.True(t, inSlice(failures, "error2: ups I did it again"))
}

func inSlice(failures []Failure, msg string) bool {
	for _, failure := range failures {
		if string(failure) == msg {
			return true
		}
	}
	return false
}
