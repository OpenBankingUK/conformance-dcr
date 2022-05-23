package compliant

const (
	expectedSSAsLen33 = 15
)

func NewDCR33(cfg DCR32Config) (Manifest, error) {
	secureClient := cfg.SecureClient
	authoriserBuilder := cfg.AuthoriserBuilder
	validator := cfg.SchemaValidator

	ssas := &cfg.SSAs
	if err := validateSSAsLen(*ssas, expectedSSAsLen33); err != nil {
		return nil, err
	}

	scenarios := Scenarios{
		DCR32ValidateOIDCConfigRegistrationURL(cfg),
		DCR32CreateSoftwareClient(cfg, secureClient, authoriserBuilder, ssas),
		DCR32DeleteSoftwareClient(cfg, secureClient, authoriserBuilder, ssas),
		DCR32CreateInvalidRegistrationRequest(cfg, secureClient, authoriserBuilder, ssas),
		DCR32RetrieveSoftwareClient(cfg, secureClient, authoriserBuilder, validator, ssas),
		DCR32RetrieveWithInvalidCredentials(cfg, secureClient, authoriserBuilder, ssas),
		DCR32UpdateSoftwareClient(cfg, secureClient, authoriserBuilder, ssas),
		DCR32UpdateSoftwareClientWithWrongId(cfg, secureClient, authoriserBuilder, ssas),
		DCR32RetrieveSoftwareClientWrongId(cfg, secureClient, authoriserBuilder, ssas),
		DCR32RegisterSoftwareWrongResponseType(cfg, secureClient, authoriserBuilder, ssas),
	}

	return NewManifest("DCR33", "1.0", scenarios)
}
