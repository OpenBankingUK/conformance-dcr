package step

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlwaysPass_Run(t *testing.T) {
	step := NewAlwaysPass()

	results := step.Run(NewContext())

	assert.True(t, results.Pass)
}
