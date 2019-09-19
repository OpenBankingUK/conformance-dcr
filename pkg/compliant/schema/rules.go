package schema

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"net/url"
	"strings"
)

const errorMessage = "The URI MUST use the https scheme; The URI MUST NOT contain a host with a value of localhost"

func isOBURLValidationRule() *validation.StringRule {
	return validation.NewStringRule(
		isOBURl,
		errorMessage,
	)
}

func isOBURl(str string) bool {
	url, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	if url.Scheme != "https" {
		return false
	}

	if url.Host == "localhost" {
		return false
	}

	if url.Host == "127.0.0.1" {
		return false
	}

	if strings.HasSuffix(url.Host, ".localhost") {
		return false
	}

	if url.Fragment != "" {
		return false
	}

	return true
}
