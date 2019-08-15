package openid

type Configuration struct {
	RegistrationEndpoint string `json:"registration_endpoint"`
	TokenEndpoint        string `json:"token_endpoint"`
}
