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

type SSAValidator struct {
	pubKeyLookup func(t *jwt.Token) (interface{}, error)
	parser       jwt.Parser
}

func NewSSAValidator(pubKeyLookup func(t *jwt.Token) (interface{}, error)) SSAValidator {
	return SSAValidator{
		pubKeyLookup: pubKeyLookup,
		parser:       jwt.Parser{ValidMethods: []string{SigningPS256, SigningES256}},
	}
}

func PublicKeyLookupFromByteSlice(pubKey []byte) func(t *jwt.Token) (interface{}, error) {
	return func(t *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(pubKey)
	}
}

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
	return SSA{SSAPayload: SSAPayload{Issuer: tkmap["iss"].(string)}}, nil
}
