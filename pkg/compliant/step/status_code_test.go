package step

import (
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"
	"github.com/stretchr/testify/assert"
)

func TestAssertStatusOk_Pass(t *testing.T) {
	ctx := context.NewContext()
	ctx.SetResponse("response", &http.Response{StatusCode: http.StatusOK})
	step := NewAssertStatus(200, "response")

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Status Code 200", result.Name)
}

func TestAssertStatusOk_FailsIfResponseNotInContext(t *testing.T) {
	ctx := context.NewContext()
	step := NewAssertStatus(200, "response")

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "getting response object from context: key not found in context", result.Message)
}

func TestAssertStatusOk_FailsIfStatusCodeIsOtherThenOk(t *testing.T) {
	ctx := context.NewContext()
	ctx.SetResponse("response", &http.Response{StatusCode: http.StatusTeapot})
	step := NewAssertStatus(200, "response")

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "status received 418", result.Message)
}
