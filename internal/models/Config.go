package models

import (
	"os"

	"gopkg.in/yaml.v3"
)

// The entire config YAML.
type Config struct {
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
		*pkgInstallConf,
		*servicesConf,
		*generalConf,
	}
}

// Grab the config from the YAML and write it to *Config.
func (conf Config) FillFromFile(filePath string) *Config {
	if filePath == "" {
		filePath = "gas.yml"
	}
	file, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		panic(err)
	}
	return &conf
}
