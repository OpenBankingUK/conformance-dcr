package schema

import (
	"fmt"
	"io"
)

type Validator interface {
	Validate(data io.Reader) []Failure
}

type Failure string

func NewValidator(version string) (Validator, error) {
	switch version {
	case "3.2":
		return responseValidator32{}, nil
	case "3.3":
		return responseValidator33{}, nil
	}
	return nil, fmt.Errorf("unknown spec version to validate schema %s", version)
}
