package openid

type Configuration struct {
	RegistrationEndpoint              string   `json:"registration_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	Issuer                            string   `json:"issuer"`
	ObjectSignAlgSupported            []string `json:"request_object_signing_alg_values_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
}
