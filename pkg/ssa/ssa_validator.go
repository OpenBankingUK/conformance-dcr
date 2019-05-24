package ssa

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
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
// the SSA header MUST comply with [RFC7519]
// signing algorithms MUST be PS256 or ES256
func NewSSAValidator(pubKeyLookup func(t *jwt.Token) (interface{}, error)) SSAValidator {
	return SSAValidator{
		pubKeyLookup: pubKeyLookup,
		parser: jwt.Parser{ValidMethods: []string{
			jwt.SigningMethodPS256.Name,
			jwt.SigningMethodES256.Name,
		}},
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
	claimMap, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return SSA{}, errors.New("unable to cast jwt.Claims to jwt.MapClaims")
	}
	if err := claimMap.Valid(); err != nil {
		return SSA{}, errors.New(fmt.Sprintf("invalid jwt claims: %v", err))
	}
	var softwareRedirectURIs []string
	for _, v := range claimMap["software_redirect_uris"].([]interface{}) {
		softwareRedirectURIs = append(softwareRedirectURIs, v.(string))
	}
	var softwareRoles []string
	for _, v := range claimMap["software_roles"].([]interface{}) {
		softwareRoles = append(softwareRoles, v.(string))
	}
	return SSA{
		// RFC7591 Header
		Typ: t.Header["typ"].(string),
		Alg: t.Header["alg"].(string),
		Kid: t.Header["kid"].(string),

		// RFC7591 Payload
		Issuer:     claimMap["iss"].(string),
		IssuedAt:   int64(claimMap["iat"].(float64)),
		JwtID:      claimMap["jti"].(string),
		SoftwareID: claimMap["software_id"].(string),

		// OB Payload
		SoftwasreEnvironment:        claimMap["software_environment"].(string),
		SoftwareMode:                claimMap["software_mode"].(string),
		SoftwareClientID:            claimMap["software_client_id"].(string),
		SoftwareClientName:          claimMap["software_client_name"].(string),
		SoftwareClientDescription:   claimMap["software_client_description"].(string),
		SoftwareClientURI:           claimMap["software_client_uri"].(string),
		SoftwareVersion:             claimMap["software_version"].(string),
		SoftwareJWKSEndpoint:        claimMap["software_jwks_endpoint"].(string),
		SoftwareJWKSRevokedEndpoint: claimMap["software_jwks_revoked_endpoint"].(string),
		SoftwareLogoURI:             claimMap["software_logo_uri"].(string),
		SoftwareOnBehalfOfOrg:       claimMap["software_on_behalf_of_org"].(string),
		SoftwarePolicyURI:           claimMap["software_policy_uri"].(string),
		SoftwareRedirectURIs:        softwareRedirectURIs,
		SoftwareRoles:               softwareRoles,
		SoftwareTermsOfServiceURI:   claimMap["software_tos_uri"].(string),
	}, nil
}
