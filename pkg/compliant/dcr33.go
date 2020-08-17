package compliant

func NewDCR33(cfg DCR32Config) (Manifest, error) {
	secureClient := cfg.SecureClient
	authoriserBuilder := cfg.AuthoriserBuilder
	validator := cfg.SchemaValidator

	scenarios := Scenarios{
		DCR32ValidateOIDCConfigRegistrationURL(cfg),
		DCR32CreateSoftwareClient(cfg, secureClient, authoriserBuilder),
		DCR32DeleteSoftwareClient(cfg, secureClient, authoriserBuilder),
		DCR32CreateInvalidRegistrationRequest(cfg, secureClient, authoriserBuilder),
		DCR32RetrieveSoftwareClient(cfg, secureClient, authoriserBuilder, validator),
		DCR32RetrieveWithInvalidCredentials(cfg, secureClient, authoriserBuilder),
		DCR32UpdateSoftwareClient(cfg, secureClient, authoriserBuilder),
		DCR32UpdateSoftwareClientWithWrongId(cfg, secureClient, authoriserBuilder),
		DCR32RetrieveSoftwareClientWrongId(cfg, secureClient, authoriserBuilder),
		DCR32RegisterSoftwareWrongResponseType(cfg, secureClient, authoriserBuilder),
	}

	return NewManifest("DCR33", "1.0", scenarios)
}
