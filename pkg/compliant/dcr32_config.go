package compliant

import (
	"crypto/rsa"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
)

type DCR32Config struct {
	OpenIDConfig openid.Configuration
	SSA          string          `json:"ssa"`
	Kid          string          `json:"kid"`
	RedirectURIs []string        `json:"redirect_uris"`
	ClientID     string          `json:"client_id"`
	PrivateKey   *rsa.PrivateKey `json:"-"`
}

func NewDCR32Config(
	openIDConfig openid.Configuration,
	ssa, kid, clientID string,
	redirectURIs []string,
	privateKey *rsa.PrivateKey,
) DCR32Config {
	return DCR32Config{
		OpenIDConfig: openIDConfig,
		SSA:          ssa,
		Kid:          kid,
		RedirectURIs: redirectURIs,
		ClientID:     clientID,
		PrivateKey:   privateKey,
	}
}
