package models

type Config struct {
	PackageInstallerConfig `yaml:",inline,omitempty"`
	ServicesConfig         `yaml:",inline,omitempty"`
	GeneralConfig          `yaml:",inline,omitempty"`
}

// This interface represents new instances of this config,
// as well as the other models we have to represent the sections
// within the YAML config file.
type IConfigs interface {
    Default()
    New()
}

func (Conf Config) Default() *Config {
	pkgInstallConf := PackageInstallerConfig{}.Default()
	servicesConf := ServicesConfig{}.Default()
	generalConf := GeneralConfig{}.Default()

	return &Config{
		*pkgInstallConf,
		*servicesConf,
		*generalConf,
	}
}
