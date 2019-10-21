package compliant

import (
	"crypto/rsa"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
)

type DCR32Config struct {
	OpenIDConfig openid.Configuration
	SSA          string
	KID          string
	RedirectURIs []string
	PrivateKey   *rsa.PrivateKey
}

func NewDCR32Config(
	openIDConfig openid.Configuration,
	ssa, kid string,
	redirectURIs []string,
	privateKey *rsa.PrivateKey,
) DCR32Config {
	return DCR32Config{
		OpenIDConfig: openIDConfig,
		SSA:          ssa,
		KID:          kid,
		RedirectURIs: redirectURIs,
		PrivateKey:   privateKey,
	}
}
