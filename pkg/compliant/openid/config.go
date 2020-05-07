package openid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Configuration struct {
	RegistrationEndpoint              *string   `json:"registration_endpoint"`
	TokenEndpoint                     string    `json:"token_endpoint"`
	RequestObjectSignAlgSupported     []string  `json:"request_object_signing_alg_values_supported"`
	TokenEndpointAuthMethodsSupported []string  `json:"token_endpoint_auth_methods_supported"`
	TokenEndpointSigningAlgSupported  *[]string `json:"token_endpoint_auth_signing_alg_values_supported"`
	ResponseTypesSupported            *[]string `json:"response_types_supported"`
}

func (c Configuration) RegistrationEndpointAsString() string {
	if c.RegistrationEndpoint == nil {
		return ""
	}

	return *c.RegistrationEndpoint
}

func Get(url string, client *http.Client) (Configuration, error) {
	r, err := client.Get(url)
	if err != nil {
		return Configuration{}, errors.Wrapf(err, "Failed to GET OpenIDConfiguration: url=%+v", url)
	}

	if r.StatusCode != http.StatusOK {
		return Configuration{}, errorFromResponse(r, url)
	}

	defer r.Body.Close()
	config := Configuration{}
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		return Configuration{}, errors.Wrap(err, "invalid OpenIDConfiguration body content")
	}

	return config, nil
}

func errorFromResponse(response *http.Response, url string) error {
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "error reading error response from GET OpenIDConfiguration")
	}

	err = response.Body.Close()
	if err != nil {
		return errors.Wrap(err, "error closing error response from GET OpenIDConfiguration")
	}

	return fmt.Errorf(
		"failed to GET OpenIDConfiguration config: url=%+v, StatusCode=%+v, body=%+v",
		url,
		response.StatusCode,
		string(responseBody),
	)
}
