package ssa

import "github.com/dgrijalva/jwt-go"

// The SSA header MUST comply with [RFC7519]
// Signing algorithms MUST be PS256 or ES256
const (
	SigningPS256 = "PS256"
	SigningES256 = "ES256"
)

// Validate gets a software statement assertion as a jwt
// parses and validates that the jwt is valid
// returns a valid SSA struct
func Validate(ssa string) (SSA, error) {
	parser := jwt.Parser{
		ValidMethods: []string{SigningPS256, SigningES256},
	}
	t, err := parser.Parse(ssa, func(t *jwt.Token) (interface{}, error) {
		return 0, nil
	})
	if err != nil {
		return SSA{}, err
	}
	return SSA{SSAPayload: SSAPayload{Issuer: t.Claims.(jwt.StandardClaims).Issuer}}, nil
}
