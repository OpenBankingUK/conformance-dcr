package config

type Config struct {
	WellknownEndpoint string `mapstructure:"wellknown_endpoint"`
	SSA               string `mapstructure:"ssa"`
}
