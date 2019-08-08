package context

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestContext_SetString(t *testing.T) {
	ctx := NewContext()
	ctx.SetString("key", "value")

	value, err := ctx.GetString("key")

	assert.NoError(t, err)
	assert.Equal(t, "value", value)
}

func TestContext_GetString_ReturnsError_IfDoesntExists(t *testing.T) {
	ctx := NewContext()
	ctx.SetString("key", "value")

	_, err := ctx.GetString("non existing key")

	assert.Equal(t, ErrKeyNotFoundInContext, err)
}

func TestContext_GetString_MultipleKeys(t *testing.T) {
	ctx := NewContext()
	ctx.SetString("key1", "value1")
	ctx.SetString("key2", "value2")

	value1, err := ctx.GetString("key1")
	value2, err := ctx.GetString("key2")

	assert.NoError(t, err)
	assert.Equal(t, "value1", value1)
	assert.Equal(t, "value2", value2)
}

func TestContext_SetInt(t *testing.T) {
	ctx := NewContext()
	ctx.SetInt("key", 42)

	value, err := ctx.GetInt("key")

	assert.NoError(t, err)
	assert.Equal(t, 42, value)
}

func TestContext_GetInt_ReturnsError_IfDoesntExists(t *testing.T) {
	ctx := NewContext()
	ctx.SetInt("key", 42)

	_, err := ctx.GetInt("non existing key")

	assert.Equal(t, ErrKeyNotFoundInContext, err)
}

func TestContext_SetResponse(t *testing.T) {
	ctx := NewContext()
	r := &http.Response{}
	ctx.SetResponse("key", r)

	value, err := ctx.GetResponse("key")

	assert.NoError(t, err)
	assert.Equal(t, r, value)
}

func TestContext_GetResponse_ReturnsError_IfDoesntExists(t *testing.T) {
	ctx := NewContext()
	r := &http.Response{}
	ctx.SetResponse("key", r)

	_, err := ctx.GetResponse("non existing key")

	assert.Equal(t, ErrKeyNotFoundInContext, err)
}
