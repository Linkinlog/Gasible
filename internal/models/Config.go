package models

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	PackageInstallerConfig `yaml:",inline,omitempty"`
	ServicesConfig         `yaml:",inline,omitempty"`
	GeneralConfig          `yaml:",inline,omitempty"`
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
