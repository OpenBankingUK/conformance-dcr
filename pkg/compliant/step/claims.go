package step

import (
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
)

type claims struct {
	stepName          string
	jwtClaimsCtxKey   string
	authoriserBuilder auth.AuthoriserBuilder
}

func NewClaims(jwtClaimsCtxKey string, authoriserBuilder auth.AuthoriserBuilder) Step {
	return claims{
		stepName:          "Generate signed software client claims",
		jwtClaimsCtxKey:   jwtClaimsCtxKey,
		authoriserBuilder: authoriserBuilder,
	}
}

func (c claims) Run(ctx Context) Result {
	debug := NewDebug()

	debug.Log("getting claims from authoriser")
	authoriser, err := c.authoriserBuilder.Build()
	if err != nil {
		return NewFailResultWithDebug(c.stepName, err.Error(), debug)
	}
	signedClaims, err := authoriser.Claims()
	if err != nil {
		return NewFailResultWithDebug(c.stepName, err.Error(), debug)
	}

	debug.Logf("setting signed claims in context var: %s", c.jwtClaimsCtxKey)
	ctx.SetString(c.jwtClaimsCtxKey, signedClaims)

	return NewPassResultWithDebug(c.stepName, debug)
}
