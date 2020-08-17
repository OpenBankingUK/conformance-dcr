package schema

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestResponseValidator33_ValidateInvalidPayload(t *testing.T) {
	validator := responseValidator33{}
	reader := bytes.NewReader([]byte(`{`))

	failures := validator.Validate(reader)

	assert.Len(t, failures, 1)
}

func TestResponseValidator33_ValidateEmpty(t *testing.T) {
	validator := responseValidator33{}
	data := []byte(`{}`)
	reader := bytes.NewReader(data)

	failures := validator.Validate(reader)

	assert.Len(t, failures, 9)
}

func TestResponseValidator33_ValidateResponse(t *testing.T) {
	validator := responseValidator33{}
	reader, err := os.Open("testdata/response33.json")
	require.NoError(t, err)

	failures := validator.Validate(reader)

	assert.Len(t, failures, 0)
}
