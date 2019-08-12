package step

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPassResult(t *testing.T) {
	passingStep := NewPassResult("some step")

	assert.True(t, passingStep.Pass)
	assert.Equal(t, "some step", passingStep.Name)
	assert.Empty(t, passingStep.Message)
}

func TestNewFailResult(t *testing.T) {
	failingTest := NewFailResult("some other step", "computer says no")

	assert.False(t, failingTest.Pass)
	assert.Equal(t, "some other step", failingTest.Name)
	assert.Equal(t, "computer says no", failingTest.Message)
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
