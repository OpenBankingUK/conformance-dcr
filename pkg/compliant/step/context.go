package step

import (
	"errors"
	"net/http"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
	dcr "github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/openid"
)

type Context interface {
	SetString(key, value string)
	GetString(key string) (string, error)
	SetInt(key string, value int)
	GetInt(key string) (int, error)
	SetResponse(key string, response *http.Response)
	GetResponse(key string) (*http.Response, error)
	SetOpenIdConfig(key string, config openid.Configuration)
	GetOpenIdConfig(key string) (openid.Configuration, error)
	SetClient(key string, client dcr.Client)
	GetClient(key string) (dcr.Client, error)
	SetGrantToken(key string, token auth.GrantToken)
	GetGrantToken(key string) (auth.GrantToken, error)
}

var ErrKeyNotFoundInContext = errors.New("key not found in context")

type context struct {
	strings       map[string]string
	ints          map[string]int
	responses     map[string]*http.Response
	openIdConfigs map[string]openid.Configuration
	clients       map[string]dcr.Client
	grantTokens   map[string]auth.GrantToken
}

func NewContext() Context {
	return &context{
		strings:       map[string]string{},
		ints:          map[string]int{},
		responses:     map[string]*http.Response{},
		openIdConfigs: map[string]openid.Configuration{},
		clients:       map[string]dcr.Client{},
		grantTokens:   map[string]auth.GrantToken{},
	}
}

func (c *context) SetString(key, value string) {
	c.strings[key] = value
}

func (c *context) GetString(key string) (string, error) {
	value, ok := c.strings[key]
	if !ok {
		return "", ErrKeyNotFoundInContext
	}
	return value, nil
}

func (c *context) SetInt(key string, value int) {
	c.ints[key] = value
}

func (c *context) GetInt(key string) (int, error) {
	value, ok := c.ints[key]
	if !ok {
		return 0, ErrKeyNotFoundInContext
	}
	return value, nil
}

func (c *context) SetResponse(key string, response *http.Response) {
	c.responses[key] = response
}

func (c *context) GetResponse(key string) (*http.Response, error) {
	value, ok := c.responses[key]
	if !ok {
		return nil, ErrKeyNotFoundInContext
	}
	return value, nil
}

func (c *context) SetOpenIdConfig(key string, config openid.Configuration) {
	c.openIdConfigs[key] = config
}

func (c *context) GetOpenIdConfig(key string) (openid.Configuration, error) {
	value, ok := c.openIdConfigs[key]
	if !ok {
		return openid.Configuration{}, ErrKeyNotFoundInContext
	}
	return value, nil
}

func (c *context) SetClient(key string, client dcr.Client) {
	c.clients[key] = client
}

func (c *context) GetClient(key string) (dcr.Client, error) {
	value, ok := c.clients[key]
	if !ok {
		return nil, ErrKeyNotFoundInContext
	}
	return value, nil
}

func (c *context) SetGrantToken(key string, token auth.GrantToken) {
	c.grantTokens[key] = token
}

func (c *context) GetGrantToken(key string) (auth.GrantToken, error) {
	value, ok := c.grantTokens[key]
	if !ok {
		return auth.GrantToken{}, ErrKeyNotFoundInContext
	}
	return value, nil
}
