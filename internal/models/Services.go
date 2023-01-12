package models

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServicesConfig struct {
	Installer  bool `yaml:"installer,omitempty"`
	Teamviewer bool `yaml:"teamviewer,omitempty"`
	Ssh        bool `yaml:"ssh,omitempty"`
	Git        bool `yaml:"git,omitempty"`
}

func (Services ServicesConfig) Default() *ServicesConfig {
	return &ServicesConfig{
		Installer:  true,
		Teamviewer: true,
		Ssh:        true,
		Git:        true,
	}
}

func (conf *ServicesConfig) Fill(filePath string) {
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
}
