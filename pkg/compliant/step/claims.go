package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
	"crypto/rsa"
	"fmt"
)

type claims struct {
	stepName           string
	jwtClaimsCtxKey    string
	openIdConfigCtxKey string
	privateKey         *rsa.PrivateKey
	ssa                string
}

func NewClaims(jwtClaimsCtxKey, openIdConfigCtxKey, ssa string, privateKey *rsa.PrivateKey) Step {
	return claims{
		stepName:           "Generate signed software client claims",
		jwtClaimsCtxKey:    jwtClaimsCtxKey,
		openIdConfigCtxKey: openIdConfigCtxKey,
		privateKey:         privateKey,
		ssa:                ssa,
	}
}

func (c claims) Run(ctx Context) Result {
	configuration, err := ctx.GetOpenIdConfig(c.openIdConfigCtxKey)
	if err != nil {
		return NewFailResult(c.stepName, fmt.Sprintf("getting openid config: %s", err.Error()))
	}

	auther := auth.NewAuthoriser(configuration, c.privateKey, c.ssa)
	signedClaims, err := auther.Claims()
	if err != nil {
		return NewFailResult(c.stepName, err.Error())
	}

	ctx.SetString(c.jwtClaimsCtxKey, signedClaims)

	return NewPassResult(c.stepName)
}
