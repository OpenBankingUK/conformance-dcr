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
	debug := NewDebug()

	debug.Logf("get openid config from ctx var: %s", c.openIdConfigCtxKey)
	configuration, err := ctx.GetOpenIdConfig(c.openIdConfigCtxKey)
	if err != nil {
		return NewFailResultWithDebug(
			c.stepName,
			fmt.Sprintf("getting openid config: %s", err.Error()),
			debug,
		)
	}

	debug.Log("creating a authoriser")
	auther := auth.NewAuthoriser(configuration, c.privateKey, c.ssa)
	signedClaims, err := auther.Claims()
	if err != nil {
		return NewFailResultWithDebug(c.stepName, err.Error(), debug)
	}

	debug.Logf("setting signed claims in context var: %s", c.jwtClaimsCtxKey)
	ctx.SetString(c.jwtClaimsCtxKey, signedClaims)

	return NewPassResultWithDebug(c.stepName, debug)
}
