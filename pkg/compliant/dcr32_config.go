package compliant

import (
	"crypto/rsa"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
)

type DCR32Config struct {
	OpenIDConfig      openid.Configuration
	SSA               string
	KID               string
	RedirectURIs      []string
	PrivateKey        *rsa.PrivateKey
	GetImplemented    bool
	PutImplemented    bool
	DeleteImplemented bool
}

func NewDCR32Config(
	openIDConfig openid.Configuration,
	ssa, kid string,
	redirectURIs []string,
	privateKey *rsa.PrivateKey,
	getImplemented bool,
	putImplemented bool,
	deleteImplemented bool,
) DCR32Config {
	return DCR32Config{
		OpenIDConfig:      openIDConfig,
		SSA:               ssa,
		KID:               kid,
		RedirectURIs:      redirectURIs,
		PrivateKey:        privateKey,
		GetImplemented:    getImplemented,
		PutImplemented:    putImplemented,
		DeleteImplemented: deleteImplemented,
	}

}
