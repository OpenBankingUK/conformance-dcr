package compliant

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

// resolves what token signing method to use based on .wellknown and supported
func responseTokenSignMethod(methods *[]string) (jwt.SigningMethod, error) {
	if methods == nil {
		return jwt.SigningMethodPS256, nil
	}

	for _, value := range *methods {
		// we only support PS256 at the moment
		if value == "PS256" {
			return jwt.SigningMethodPS256, nil
		}
	}

	return nil, errors.New("PS256 token sign method not found")
}
