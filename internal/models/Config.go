package models

import (
	"os"

	"gopkg.in/yaml.v3"
)

// The entire config YAML.
type Config struct {
	FilePath               string
	PackageInstallerConfig `yaml:",inline,omitempty"`
	ServicesConfig         `yaml:",inline,omitempty"`
	GeneralConfig          `yaml:",inline,omitempty"`
}

// Create the defaults and write them to *Config.
func (Conf Config) Default() *Config {
	pkgInstallConf := PackageInstallerConfig{}.Default()
	servicesConf := ServicesConfig{}.Default()
	generalConf := GeneralConfig{}.Default()

	return &Config{
		"gas.yml",
		*pkgInstallConf,
		*servicesConf,
		*generalConf,
	}
}

// Grab the config from the YAML and write it to the given struct.
func (conf *Config)FillFromFile() {
    conf.FilePath = "gas.yml"
	file, err := os.ReadFile(conf.FilePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		panic(err)
	}
}

