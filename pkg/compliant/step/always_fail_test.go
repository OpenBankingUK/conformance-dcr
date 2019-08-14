package step

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlwaysFail_Run(t *testing.T) {
	step := NewAlwaysFail()

	results := step.Run(NewContext())

	assert.False(t, results.Pass)
}
