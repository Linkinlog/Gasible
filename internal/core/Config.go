package core

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// The entire config YAML.
type ConfigModel struct {
	PackageManagerConfig `yaml:",inline,omitempty"`
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
func NewConfigWithDefaults() ConfigModel {
	pkgInstallConf := PackageManagerConfig{}.Default()
	servicesConf := ServicesConfig{}.Default()
	generalConf := GeneralConfig{}.Default()
	globalOpts := &GlobalOpts{}

	return ConfigModel{
		*pkgInstallConf,
		*servicesConf,
		*generalConf,
		*globalOpts,
	}
}

// Create the defaults and write them to *Config.
//
// TODO better handling of this.
func NewConfigFromFile() *ConfigModel {
	conf := &ConfigModel{
		PackageManagerConfig{},
		ServicesConfig{},
		GeneralConfig{},
		GlobalOpts{},
	}
	if err := conf.LoadFromFile(); err != nil {
		panic(err)
	}
	return conf
}

// This is the ConfigModel for the running application.
var Config = NewConfigFromFile()

// Grab the config from the YAML and write it to the given struct.
func (conf *ConfigModel) LoadFromFile() error {
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
