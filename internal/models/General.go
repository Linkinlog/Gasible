package models

import (
	"os"

	"gopkg.in/yaml.v3"
)

type TVCreds struct {
	User string `yaml:"user,omitempty"`
	Pass string `yaml:"pass,omitempty"`
}

type GeneralConfig struct {
	Hostname        string  `yaml:"hostname,omitempty"`
	IP              string  `yaml:"staticIP,omitempty"`
	Mask            string  `yaml:"mask,omitempty"`
	TeamViewerCreds TVCreds `yaml:"TeamViewerCreds,omitempty"`
}

func (General GeneralConfig) Default() *GeneralConfig {
	return &GeneralConfig{
		Hostname: "development-station",
		IP:       "192.168.4.20",
		Mask:     "255.255.255.0",
		TeamViewerCreds: TVCreds{
			User: "username",
			Pass: "password",
		},
	}
}

func (conf *GeneralConfig) Fill(filePath string) {
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
