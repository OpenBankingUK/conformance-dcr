package step

import (
	"errors"
	"net/http"
)

type Context interface {
	SetString(key, value string)
	GetString(key string) (string, error)
	SetInt(key string, value int)
	GetInt(key string) (int, error)
	SetResponse(key string, response *http.Response)
	GetResponse(key string) (*http.Response, error)
}

var ErrKeyNotFoundInContext = errors.New("key not found in context")

type context struct {
	strings   map[string]string
	ints      map[string]int
	responses map[string]*http.Response
}

func NewContext() Context {
	return &context{
		strings:   map[string]string{},
		ints:      map[string]int{},
		responses: map[string]*http.Response{},
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
