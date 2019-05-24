package ssa

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
			jwt.SigningMethodPS256.Alg(),
			jwt.SigningMethodES256.Alg(),
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
		jwkKid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, errors.New("unable to cast jwk kid to string")
		}
		res, err := http.Get(jwkEndpoint)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		var jwk map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&jwk)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("unable to parse json from jwk endpoint %s err: %v", jwkEndpoint, err))
		}
		for _, v := range jwk["keys"].([]interface{}) {
			v, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if v["kid"].(string) != jwkKid {
				continue
			}
			certURI, ok := v["x5u"].(string)
			if !ok {
				return nil, errors.New("unable to cast `x5u` parameter to string")
			}
			res, err := http.Get(certURI)
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()
			certBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("unable to download certificate: %v", err))
			}
			return jwt.ParseRSAPublicKeyFromPEM(certBytes)
		}
		return nil, errors.New(fmt.Sprintf("unable to find key with kid %s in jwks endpoint key store %s", jwkKid, jwkEndpoint))
	}
}

// Validate gets a software statement assertion as a jwt
// parses and validates that the jwt is valid
// returns a valid SSA struct
func (v SSAValidator) Validate(ssa string) (SSA, error) {
	t, err := v.parser.ParseWithClaims(ssa, &SSA{}, v.pubKeyLookup)
	if err != nil {
		return SSA{}, err
	}
	if claims, ok := t.Claims.(*SSA); ok && t.Valid {
		claims.Typ = t.Header["typ"].(string)
		claims.Alg = t.Header["alg"].(string)
		claims.Kid = t.Header["kid"].(string)

		return *claims, nil
	} else {
		return SSA{}, err
	}
}
