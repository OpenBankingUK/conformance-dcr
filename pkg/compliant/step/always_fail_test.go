package step

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"
	"github.com/stretchr/testify/assert"
)

func TestAlwaysFail_Run(t *testing.T) {
	step := NewAlwaysFail()

	results := step.Run(context.NewContext())

	assert.False(t, results.Pass)
}
