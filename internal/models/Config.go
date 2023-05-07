package models

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// The entire config YAML.
type Config struct {
	PackageInstallerConfig `yaml:",inline,omitempty"`
	ServicesConfig         `yaml:",inline,omitempty"`
	GeneralConfig          `yaml:",inline,omitempty"`
	GlobalOpts             `yaml:"-"`
}

// Options to embed on our config that
// we may need throughout execution.
type GlobalOpts struct {
	FilePath string
	NoOp     bool
}

// Create the defaults and write them to *Config.
func NewConfigWithDefaults() *Config {
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
func (conf *Config) LoadFromFile() error {
	logFile := conf.GlobalOpts.FilePath
	_, err := os.Stat(logFile)
	if err != nil {
		return errors.New("Config file found.")
	}
	file, err := os.ReadFile(logFile)
	if err != nil {
		return errors.New("Could read from config file.")
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return errors.New("Could not Unmarshal Config file.")
	}

	return nil
}
