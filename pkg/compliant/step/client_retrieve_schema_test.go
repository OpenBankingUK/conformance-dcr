package step

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/schema"
	"github.com/stretchr/testify/assert"
)

func TestNewClientRetrieveSchema(t *testing.T) {
	validator := &stubValidator{}
	ctx := NewContext()
	body := ioutil.NopCloser(strings.NewReader(`{}`))
	ctx.SetResponse("responseCtxKey", &http.Response{Body: body})
	step := NewClientRetrieveSchema("responseCtxKey", validator)

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Validate client response schema", result.Name)
	assert.Equal(t, "", result.FailReason)
	assert.True(t, validator.called)
}

func TestNewClientRetrieveSchema_FailsMissingCtxResponse(t *testing.T) {
	validator := &stubValidator{}
	ctx := NewContext()
	step := NewClientRetrieveSchema("responseCtxKey", validator)

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "getting response object from context: key not found in context", result.FailReason)
}

func TestNewClientRetrieveSchema_MapsErrors(t *testing.T) {
	validator := &stubValidator{failures: []schema.Failure{"ups"}}
	ctx := NewContext()
	body := ioutil.NopCloser(strings.NewReader(`{}`))
	ctx.SetResponse("responseCtxKey", &http.Response{Body: body})
	step := NewClientRetrieveSchema("responseCtxKey", validator)

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "schema invalid: ups, ups", result.FailReason)
}

type stubValidator struct {
	called   bool
	failures []schema.Failure
}

func (s *stubValidator) Validate(reader io.Reader) []schema.Failure {
	s.called = true
	return s.failures
}
