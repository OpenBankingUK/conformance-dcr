package compliant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTester(t *testing.T) {
	tester := NewTester("", false)
	assert.IsType(t, colourTester{}, tester)

	tester = NewTester("something", false)
	assert.IsType(t, filteredTester{}, tester)
}
