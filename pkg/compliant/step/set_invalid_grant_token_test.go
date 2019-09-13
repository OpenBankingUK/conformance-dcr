package step

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetInvalidGrantToken_Run(t *testing.T) {
	ctx := NewContext()
	step := NewSetInvalidGrantToken("grantTokenCtxKey")

	result := step.Run(ctx)

	token, err := ctx.GetGrantToken("grantTokenCtxKey")
	require.NoError(t, err)
	assert.True(t, result.Pass)
	assert.Equal(t, "Set invalid grant token", result.Name)
	assert.Equal(t, "", token.AccessToken)
	assert.Equal(t, "", token.TokenType)
}
