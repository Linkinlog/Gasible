// Config model file.
// This file defines the Config model and its methods.
// The config model is a representation of the whole YAML configuration.
package models

import (
	"errors"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// The entire config YAML.
type Config struct {
	PackageInstallerConfig `yaml:",inline,omitempty"`
	GitServiceConfig       `yaml:"Git Config,omitempty"`
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
func (Conf Config) Default() *Config {
	pkgInstallConf := PackageInstallerConfig{}.Default()
	gitServiceConf := GitServiceConfig{}.Default()
	servicesConf := ServicesConfig{}.Default()
	generalConf := GeneralConfig{}.Default()
	globalOpts := &GlobalOpts{}

	return &Config{
		*pkgInstallConf,
		*gitServiceConf,
		*servicesConf,
		*generalConf,
		*globalOpts,
	}
}

// Grab the config from the YAML and write it to the given struct.
func (conf *Config) FillFromFile() error {
	yamlFile := conf.GlobalOpts.FilePath
	_, err := os.Stat(yamlFile)
	if err != nil {
		return errors.New("Panic: No file " + yamlFile + " found.\nRun the generate command to make a new config")
	}
	file, err := os.ReadFile(yamlFile)
	if err != nil {
		log.Fatalf("\nPanic: Could read file %s", yamlFile)
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return err
	}
	return nil
}
