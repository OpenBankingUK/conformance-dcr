package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClientRegister(t *testing.T) {
	// creating a stub server that expects a JWT body posted
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.Equal(t, req.URL.String(), "/some/path")

		// is it a JWT body?
		require.Equal(t, "application/jwt", req.Header.Get("Content-Type"))

		// does it have the JWT body?
		body, err := ioutil.ReadAll(req.Body)
		require.NoError(t, err)
		assert.Equal(t, "ssa", string(body))

		_, err = rw.Write([]byte(`OK`))
		require.NoError(t, err)
	}))
	defer server.Close()

	ctx := NewContext()
	url := server.URL + "/some/path"
	ctx.SetOpenIdConfig("openIdConfigKey", openid.Configuration{
		RegistrationEndpoint: url,
		TokenEndpoint:        "",
	})
	step := NewClientRegister("openIdConfigKey", "ssa", "responseKey", server.Client())

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Software client register", result.Name)
	assert.Equal(t, "", result.Message)

	// assert that response in now in ctx
	_, err := ctx.GetResponse("responseKey")
	assert.NoError(t, err)
}

func TestNewClientRegister_HandlesHttpErrors(t *testing.T) {
	ctx := NewContext()
	ctx.SetOpenIdConfig("openIdConfigKey", openid.Configuration{
		RegistrationEndpoint: "invalid url",
		TokenEndpoint:        "",
	})
	step := NewClientRegister("openIdConfigKey", "ssa", "responseKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "making jwt post request: Post invalid%20url: unsupported protocol scheme \"\"", result.Message)
}

func TestNewClientRegister_HandlesCreateRequestError(t *testing.T) {
	ctx := NewContext()
	ctx.SetOpenIdConfig("openIdConfigKey", openid.Configuration{
		RegistrationEndpoint: string(0x7f),
		TokenEndpoint:        "",
	})
	step := NewClientRegister("openIdConfigKey", "ssa", "responseKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "creating jwt post request: parse \u007f: net/url: invalid control character in URL", result.Message)
}

func TestNewClientRegister_HandlesOpenIdConfigNotInContext(t *testing.T) {
	ctx := NewContext()
	step := NewClientRegister("openIdConfigKey", "ssa", "responseKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "getting openid config: key not found in context", result.Message)
}
