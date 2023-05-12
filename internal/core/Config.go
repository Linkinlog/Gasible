package core

import (
	"errors"

	"gopkg.in/yaml.v3"
)

// The entire config YAML.
type ConfigModel struct {
	PackageManagerConfig `yaml:",inline,omitempty"`
	ServicesConfig       `yaml:",inline,omitempty"`
	GeneralConfig        `yaml:",inline,omitempty"`
	GlobalOpts           `yaml:"-"`
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
func NewConfigFromFile(configFile string) *ConfigModel {
	conf := &ConfigModel{
		PackageManagerConfig{},
		ServicesConfig{},
		GeneralConfig{},
		GlobalOpts{},
	}
	if err := conf.LoadFromFile(configFile); err != nil {
		// TODO not a fan of panicing but itll work for now
		panic(err)
	}
	return conf
}

// Grab the config from the YAML and write it to the given struct.
func (conf *ConfigModel) LoadFromFile(configFile string) error {
	err := yaml.Unmarshal([]byte(configFile), &conf)
	if err != nil {
		return errors.New("Could not Unmarshal Config file.")
	}

	return nil
}
