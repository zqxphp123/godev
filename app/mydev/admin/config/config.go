package config

import (
	"mydev/app/pkg/options"
	"mydev/pkg/app"
	cliflag "mydev/pkg/common/cli/flag"
	"mydev/pkg/log"
)

type Config struct {
	Log      *log.Options             `json:"log" mapstructure:"log"`
	Server   *options.ServerOptions   `json:"server" mapstructure:"server"`
	Registry *options.RegistryOptions `json:"registry" mapstructure:"registry"`
}

func (c *Config) Validate() []error {
	var errors []error
	errors = append(errors, c.Log.Validate()...)
	errors = append(errors, c.Server.Validate()...)
	errors = append(errors, c.Registry.Validate()...)
	return errors
}

func (c *Config) Flags() (fss cliflag.NamedFlagSets) {
	c.Log.AddFlags(fss.FlagSet("logs"))
	c.Server.AddFlags(fss.FlagSet("server"))
	c.Registry.AddFlags(fss.FlagSet("registry"))
	return fss
}
func New() *Config {
	return &Config{
		Log:      log.NewOptions(),
		Server:   options.NewServerOptions(),
		Registry: options.NewRegistryOptions(),
	}
}

var _ app.CliOptions = &Config{}
