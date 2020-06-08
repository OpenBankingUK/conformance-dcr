package compliant

import "github.com/pkg/errors"

func IsSupportedSpecVersion(version string) bool {
	return version == "3.2" || version == "3.3"
}

func NewSpecManifest(version string, cfg DCR32Config) (Manifest, error) {
	switch version {
	case "3.2":
		return NewDCR32(cfg)
	case "3.3":
		return nil, errors.New("not implemented yet")
	}
	return nil, errors.New("specification version  not supported")
}
