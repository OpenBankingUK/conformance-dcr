package step

import (
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
)

type setInvalidGrantToken struct {
	stepName         string
	grantTokenCtxKey string
}

func NewSetInvalidGrantToken(grantTokenCtxKey string) Step {
	return setInvalidGrantToken{
		stepName:         "Set invalid grant token",
		grantTokenCtxKey: grantTokenCtxKey,
	}
}

func (s setInvalidGrantToken) Run(ctx Context) Result {
	debug := NewDebug()

	token := auth.GrantToken{
		AccessToken: "",
		TokenType:   "",
		ExpiresIn:   0,
	}

	debug.Logf("setting invalid client credentials token in context var: %s", s.grantTokenCtxKey)
	ctx.SetGrantToken(s.grantTokenCtxKey, token)

	return NewPassResultWithDebug(s.stepName, debug)
}
