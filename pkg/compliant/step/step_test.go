package step

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPassResult(t *testing.T) {
	passingStep := NewPassResult("some step")

	assert.True(t, passingStep.Pass)
	assert.Equal(t, "some step", passingStep.Name)
	assert.Empty(t, passingStep.FailReason)
}

func TestNewFailResult(t *testing.T) {
	failingTest := NewFailResult("some other step", "computer says no")

	assert.False(t, failingTest.Pass)
	assert.Equal(t, "some other step", failingTest.Name)
	assert.Equal(t, "computer says no", failingTest.FailReason)
}

func TestResults_Fail_False_All_Passing(t *testing.T) {
	passingSteps := Results{
		NewPassResult("some step"),
		NewPassResult("and some more"),
	}

	assert.False(t, passingSteps.Fail())
}

func TestResults_Fail_One_Failing(t *testing.T) {
	passingSteps := Results{
		NewPassResult("some step"),
		NewFailResult("this one failed", "failed for some reason"),
		NewPassResult("and some more"),
	}

	assert.True(t, passingSteps.Fail())
}

func TestDebugMessages_Log(t *testing.T) {
	debug := NewDebug()

	debug.Log("It all starts here.")
	debug.Log("The End!")

	assert.Len(t, debug.Item, 2)
}

func TestDebugMessages_Logf(t *testing.T) {
	debug := NewDebug()

	debug.Logf("What's %s then %s?", "better", "beer")

	assert.Len(t, debug.Item, 1)
	assert.Equal(t, "What's better then beer?", debug.Item[0].Message)
}
