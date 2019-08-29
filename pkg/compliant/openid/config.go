package openid

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type Configuration struct {
	RegistrationEndpoint              string   `json:"registration_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	Issuer                            string   `json:"issuer"`
	ObjectSignAlgSupported            []string `json:"request_object_signing_alg_values_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
}

func Get(url string, client *http.Client) (Configuration, error) {
	resp, err := client.Get(url)
	if err != nil {
		return Configuration{}, errors.Wrapf(err, "Failed to GET OpenIDConfiguration: url=%+v", url)
	}

	if resp.StatusCode != http.StatusOK {
		responseBody, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return Configuration{}, errors.Wrap(err, "error reading error response from GET OpenIDConfiguration")
		}

		return Configuration{}, fmt.Errorf(
			"failed to GET OpenIDConfiguration config: url=%+v, StatusCode=%+v, body=%+v",
			url,
			resp.StatusCode,
			string(responseBody),
		)
	}

	defer resp.Body.Close()
	config := Configuration{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return config, errors.Wrap(err, "invalid OpenIDConfiguration body content")
	}
	return config, nil
}
