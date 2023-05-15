package core

import (
	"errors"

	"gopkg.in/yaml.v3"
)

// ConfigModel The entire config YAML.
type ConfigModel struct {
	PackageManagerConfig `yaml:",inline,omitempty"`
	ServicesConfig       `yaml:",inline,omitempty"`
	GeneralConfig        `yaml:",inline,omitempty"`
	GlobalOpts           `yaml:"-"`
}

// GlobalOpts Options to embed on our config that
// we may need throughout execution.
type GlobalOpts struct {
	FilePath string
	NoOp     bool
}

var CurrentConfig = ConfigModel{}

func SetConfig(model *ConfigModel) {
	CurrentConfig = *model
}

func GetConfig() *ConfigModel {
	return &CurrentConfig
}

// NewConfigWithDefaults creates the defaults and write them to *Config.
func NewConfigWithDefaults() ConfigModel {
	pkgInstallConf := PackageManagerConfig{}
	servicesConf := ServicesConfig{}
	generalConf := GeneralConfig{}
	globalOpts := &GlobalOpts{}

	return ConfigModel{
		*pkgInstallConf.Default(),
		*servicesConf.Default(),
		*generalConf.Default(),
		*globalOpts,
	}
}

// NewConfigFromFile creates the defaults and write them to *Config.
func NewConfigFromFile(configFile string) *ConfigModel {
	conf := &ConfigModel{
		PackageManagerConfig{},
		ServicesConfig{},
		GeneralConfig{},
		GlobalOpts{},
	}
	if err := conf.LoadFromFile(configFile); err != nil {
		// TODO not a fan of panicking but it'll work for now
		panic(err)
	}
	return conf
}

// LoadFromFile grabs the config from the YAML and write it to the given struct.
func (conf *ConfigModel) LoadFromFile(configFile string) error {
	err := yaml.Unmarshal([]byte(configFile), &conf)
	if err != nil {
		return errors.New("could not Unmarshal Config file")
	}

	return nil
}
