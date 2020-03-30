package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
)

func TestDrainBody(t *testing.T) {
	reader := ioutil.NopCloser(strings.NewReader("hello"))

	copyOne, copyTwo, err := DrainBody(reader)

	require.NoError(t, err)
	dataOne, err := ioutil.ReadAll(copyOne)
	assert.Equal(t, []byte("hello"), dataOne)
	require.NoError(t, err)
	dataTwo, err := ioutil.ReadAll(copyTwo)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), dataTwo)
}
