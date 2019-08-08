package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlwaysPass_Run(t *testing.T) {
	step := NewAlwaysPass(1)

	result := step.Run(context.NewContext())

	assert.True(t, result.Pass)
	assert.Equal(t, 1, step.Order())
}
