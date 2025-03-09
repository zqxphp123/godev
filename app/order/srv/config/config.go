package config

import (
	"encoding/json"
	"mydev/app/pkg/options"

	cliflag "mydev/pkg/common/cli/flag"
	"mydev/pkg/log"
)

type Config struct {
	MySQLOptions *options.MySQLOptions     `json:"mysql"     mapstructure:"mysql"`
	Log          *log.Options              `json:"log"     mapstructure:"log"`
	Server       *options.ServerOptions    `json:"server"     mapstructure:"server"`
	Telemetry    *options.TelemetryOptions `json:"telemetry" mapstructure:"telemetry"`
	Registry     *options.RegistryOptions  `json:"consul" mapstructure:"consul"`
	Dtm          *options.DtmOptions       `json:"dtm" mapstructure:"dtm"`
}

func New() *Config {
	//配置默认初始化
	return &Config{
		MySQLOptions: options.NewMySQLOptions(),
		Log:          log.NewOptions(),
		Server:       options.NewServerOptions(),
		Telemetry:    options.NewTelemetryOptions(),
		Registry:     options.NewRegistryOptions(),
		Dtm:          options.NewDtmOptionsOptions(),
	}
}

// Flags returns flags for a specific APIServer by section name.
func (o *Config) Flags() (fss cliflag.NamedFlagSets) {
	o.Server.AddFlags(fss.FlagSet("server"))
	o.Log.AddFlags(fss.FlagSet("logs"))
	o.Telemetry.AddFlags(fss.FlagSet("telemetry"))
	o.Registry.AddFlags(fss.FlagSet("registry"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.Dtm.AddFlags(fss.FlagSet("dtm"))
	return fss
}

func (o *Config) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

func (o *Config) Validate() []error {
	var errs []error

	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.Server.Validate()...)
	errs = append(errs, o.Telemetry.Validate()...)
	errs = append(errs, o.Registry.Validate()...)
	errs = append(errs, o.Dtm.Validate()...)
	return errs
}
