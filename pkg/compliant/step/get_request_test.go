package step

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRequest_Pass(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.Equal(t, req.URL.String(), "/some/path")
		_, err := rw.Write([]byte(`OK`))
		require.NoError(t, err)
	}))
	defer server.Close()
	ctx := context.NewContext()
	url := server.URL + "/some/path"
	step := NewGetRequest(url, "response", server.Client())

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "GET request "+url, result.Name)
	assert.Equal(t, "", result.Message)
}

func TestGetRequest_FailsIfHttpCallFails(t *testing.T) {
	ctx := context.NewContext()
	step := NewGetRequest("invalid_url", "response", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "Get invalid_url: unsupported protocol scheme \"\"", result.Message)
}

func TestGetRequest_SetsResponseInContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.Equal(t, req.URL.String(), "/some/path")
		_, err := rw.Write([]byte(`OK`))
		require.NoError(t, err)
	}))
	defer server.Close()
	ctx := context.NewContext()
	url := server.URL + "/some/path"
	step := NewGetRequest(url, "response", server.Client())
	step.Run(ctx)

	r, err := ctx.GetResponse("response")

	require.NoError(t, err)
	body, err := ioutil.ReadAll(r.Body)
	require.NoError(t, err)
	assert.Equal(t, []byte(`OK`), body)
}
