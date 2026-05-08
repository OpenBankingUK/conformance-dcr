package http

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDrainBody(t *testing.T) {
	reader := io.NopCloser(strings.NewReader("hello"))

	copyOne, copyTwo, err := DrainBody(reader)

	require.NoError(t, err)
	dataOne, err := io.ReadAll(copyOne)
	assert.Equal(t, []byte("hello"), dataOne)
	require.NoError(t, err)
	dataTwo, err := io.ReadAll(copyTwo)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), dataTwo)
}
