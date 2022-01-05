package step

import (
	"net/http"
	"testing"

	dcr "github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/openid"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
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

func TestContext_SetOpenIdConfig(t *testing.T) {
	ctx := NewContext()
	config := openid.Configuration{}
	ctx.SetOpenIdConfig("key", config)

	value, err := ctx.GetOpenIdConfig("key")

	assert.NoError(t, err)
	assert.Equal(t, config, value)
}

func TestContext_GetOpenIdConfig_ReturnsError_IfDoesntExists(t *testing.T) {
	ctx := NewContext()
	config := openid.Configuration{}
	ctx.SetOpenIdConfig("key", config)

	_, err := ctx.GetOpenIdConfig("non existing key")

	assert.Equal(t, ErrKeyNotFoundInContext, err)
}

func TestContext_SetClient(t *testing.T) {
	ctx := NewContext()
	client := dcr.NewClientSecretBasic("id", "tokenEndpoint", "Token")
	ctx.SetClient("key", client)

	value, err := ctx.GetClient("key")

	assert.NoError(t, err)
	assert.Equal(t, client, value)
}

func TestContext_GetClient_ReturnsError_IfDoesntExists(t *testing.T) {
	ctx := NewContext()
	client := dcr.NewClientSecretBasic("id", "tokenEndpoint", "Token")
	ctx.SetClient("key", client)

	_, err := ctx.GetClient("non existing key")

	assert.Equal(t, ErrKeyNotFoundInContext, err)
}
