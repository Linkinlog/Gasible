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
	GlobalOpts
}

// Options to embed on our config that
// we may need throughout execution.
type GlobalOpts struct {
	FilePath string
	NoOp     bool
}

// Create the defaults and write them to *Config.
func (Conf Config) Default() *Config {
	pkgInstallConf := PackageInstallerConfig{}.Default()
	servicesConf := ServicesConfig{}.Default()
	generalConf := GeneralConfig{}.Default()
	globalOpts := &GlobalOpts{}

	return &Config{
		*pkgInstallConf,
		*servicesConf,
		*generalConf,
		*globalOpts,
	}
}

// Grab the config from the YAML and write it to the given struct.
func (conf *Config) FillFromFile() {
	file, err := os.ReadFile(conf.GlobalOpts.FilePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		panic(err)
	}
}
