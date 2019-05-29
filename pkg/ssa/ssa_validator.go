package ssa

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
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
func PublicKeyLookupFromJWKSEndpoint(client *http.Client) func(t *jwt.Token) (interface{}, error) {
	return func(t *jwt.Token) (interface{}, error) {
		tkmap, ok := t.Claims.(*SSA)
		if !ok {
			return nil, errors.New("unable to cast token claim to map[string]interface{}")
		}
		jwkEndpoint := tkmap.SoftwareJWKSEndpoint
		jwkKid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, errors.New("unable to cast jwk kid to string")
		}
		res, err := client.Get(jwkEndpoint)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to retrieve data from jwks endpoint %s", jwkEndpoint)
		}
		defer res.Body.Close()
		var jwk map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&jwk)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse json from jwk endpoint %s", jwkEndpoint)
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
			res, err := client.Get(certURI)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to download certificate from URI %s", certURI)
			}
			defer res.Body.Close()
			certBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to read certificate bytes after download")
			}
			return jwt.ParseRSAPublicKeyFromPEM(certBytes)
		}
		return nil, errors.Errorf("unable to find key with kid %s in jwks endpoint key store %s. Got response %v", jwkKid, jwkEndpoint, jwk)
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
	if ssa, ok := t.Claims.(*SSA); ok && t.Valid {
		ssa.Typ = t.Header["typ"].(string)
		ssa.Alg = t.Header["alg"].(string)
		ssa.Kid = t.Header["kid"].(string)

		if err := ssa.Validate(); err != nil {
			return SSA{}, errors.Wrap(err, "ssa is not valid")
		}

		return *ssa, nil
	} else {
		return SSA{}, err
	}
}
