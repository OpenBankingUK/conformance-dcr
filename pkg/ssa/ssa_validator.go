package ssa

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// The SSA header MUST comply with [RFC7519]
// Signing algorithms MUST be PS256 or ES256
const (
	SigningPS256 = "PS256"
	SigningES256 = "ES256"
)

// SSAValidator is a struct responsible for verification
// parsing and decryption of the SSA jwt
// it can be initialised with a custom publicKeyLookup function
// it can be initialised with a custom instance of jwt.Parser
type SSAValidator struct {
	pubKeyLookup func(t *jwt.Token) (interface{}, error)
	parser       jwt.Parser
}

// NewSSAValidator returns a new instance of SSAValidator with a specified
// pub Key lookup function that can be passed as parameter
// the constructor also defines the allowed valid methods to verify the jwt
func NewSSAValidator(pubKeyLookup func(t *jwt.Token) (interface{}, error)) SSAValidator {
	return SSAValidator{
		pubKeyLookup: pubKeyLookup,
		parser:       jwt.Parser{ValidMethods: []string{SigningPS256, SigningES256}},
	}
}

// PublicKeyLookupFromByteSlice returns a function which returns
// the same public key it has got as parameter
// mostly used for testing and debugging purposes
func PublicKeyLookupFromByteSlice(pubKey []byte) func(t *jwt.Token) (interface{}, error) {
	return func(t *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(pubKey)
	}
}

// PublicKeyLookupFromJWKSEndpoint returns a function which looks up the public key
// from a jwk endpoint specified in the jwt token
// it uses the kid to retrieve the right public key to verify the validity of the jwt
func PublicKeyLookupFromJWKSEndpoint() func(t *jwt.Token) (interface{}, error) {
	return func(t *jwt.Token) (interface{}, error) {
		tkmap, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("unable to cast token claim to map[string]interface{}")
		}
		jwkEndpoint, ok := tkmap["software_jwks_endpoint"].(string)
		if !ok {
			return nil, errors.New("unable to cast jwk endpoint to string")
		}
		res, err := http.Get(jwkEndpoint)
		if err != nil {
			return nil, err
		}
		var jwk map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&jwk)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("unable to parse json from jwk endpoint %s err: %v", jwkEndpoint, err))
		}
		// todo
		// lookup for x5u with kid == tkmap["kid"]
		return jwt.ParseRSAPublicKeyFromPEM([]byte(``))
	}
}

// Validate gets a software statement assertion as a jwt
// parses and validates that the jwt is valid
// returns a valid SSA struct
func (v SSAValidator) Validate(ssa string) (SSA, error) {
	t, err := v.parser.Parse(ssa, v.pubKeyLookup)
	if err != nil {
		return SSA{}, err
	}
	tkmap, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return SSA{}, errors.New("unable to cast jwt.Claims to jwt.MapClaims")
	}
	return SSA{Issuer: tkmap["iss"].(string)}, nil
}
