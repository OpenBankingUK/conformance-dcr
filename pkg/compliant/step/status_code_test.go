package step

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertStatusOk_Pass(t *testing.T) {
	ctx := NewContext()
	ctx.SetResponse("response", &http.Response{StatusCode: http.StatusOK})
	step := NewAssertStatus(http.StatusOK, "response")

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Assert status code 200", result.Name)
}

func TestAssertStatusOk_FailsIfResponseNotInContext(t *testing.T) {
	ctx := NewContext()
	step := NewAssertStatus(http.StatusOK, "response")

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "getting response object from context: key not found in context", result.Message)
}

func TestAssertStatusOk_FailsIfStatusCodeIsOtherThenOk(t *testing.T) {
	ctx := NewContext()
	ctx.SetResponse("response", &http.Response{StatusCode: http.StatusTeapot})
	step := NewAssertStatus(http.StatusOK, "response")

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "status received 418", result.Message)
}
