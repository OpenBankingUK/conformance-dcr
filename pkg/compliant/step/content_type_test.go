package step

import (
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"
	"github.com/stretchr/testify/assert"
)

func TestAssertContentType_Pass(t *testing.T) {
	ctx := context.NewContext()
	headers := http.Header{"Content-Type": []string{"application/vorgon"}}
	ctx.SetResponse("response", &http.Response{Header: headers})
	step := NewAssertContentType("response", "application/vorgon")

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Content-Type header is application/vorgon", result.Name)
}

func TestAssertContentType_FailsIfResponseNotInContext(t *testing.T) {
	ctx := context.NewContext()
	step := NewAssertContentType("response", "application/vorgon")

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "getting response object from context: key not found in context", result.Message)
}

func TestAssertContentType_FailsIfHeaderIsNotInResponse(t *testing.T) {
	ctx := context.NewContext()
	ctx.SetResponse("response", &http.Response{})
	step := NewAssertContentType("response", "application/vorgon")

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "Content-Type header is not present", result.Message)
}

func TestAssertContentType_FailsIfStatusCodeIsOtherThenOk(t *testing.T) {
	ctx := context.NewContext()
	headers := http.Header{"Content-Type": []string{"application/klingon"}}
	ctx.SetResponse("response", &http.Response{Header: headers})
	step := NewAssertContentType("response", "application/vorgon")

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "Content-Type is 'application/klingon'", result.Message)
}
