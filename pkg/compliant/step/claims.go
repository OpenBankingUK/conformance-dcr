package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
)

type claims struct {
	stepName        string
	jwtClaimsCtxKey string
	authoriser      auth.Authoriser
}

func NewClaims(jwtClaimsCtxKey string, authoriser auth.Authoriser) Step {
	return claims{
		stepName:        "Generate signed software client claims",
		jwtClaimsCtxKey: jwtClaimsCtxKey,
		authoriser:      authoriser,
	}
}

func (c claims) Run(ctx Context) Result {
	debug := NewDebug()

	debug.Log("getting claims from authoriser")
	signedClaims, err := c.authoriser.Claims()
	if err != nil {
		return NewFailResultWithDebug(c.stepName, err.Error(), debug)
	}

	debug.Logf("setting signed claims in context var: %s", c.jwtClaimsCtxKey)
	ctx.SetString(c.jwtClaimsCtxKey, signedClaims)

	return NewPassResultWithDebug(c.stepName, debug)
}
