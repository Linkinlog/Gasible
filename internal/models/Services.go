package models

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Information relatie to the services we will configure.
type ServicesConfig struct {
	Installer  bool `yaml:"installer,omitempty"`
	Teamviewer bool `yaml:"teamviewer,omitempty"`
	Ssh        bool `yaml:"ssh,omitempty"`
	Git        bool `yaml:"git,omitempty"`
}

// Sets the defaults.
func (ServicesConfig) Default() *ServicesConfig {
	return &ServicesConfig{
		Installer:  true,
		Teamviewer: true,
		Ssh:        true,
		Git:        true,
	}
}

// Grabs config from YAML and fills the struct with it.
func (conf ServicesConfig) FillFromFile(filePath string) *ServicesConfig {
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
