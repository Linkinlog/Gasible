package models

import (
	"os"

	"gopkg.in/yaml.v3"
)

type PackageInstallerConfig struct {
	Manager  string   `yaml:"pkg-manager,omitempty"`
	Packages []string `yaml:"packages,omitempty"`
}

func (pkgInstallConf PackageInstallerConfig) Default() *PackageInstallerConfig{
	return &PackageInstallerConfig{
		Manager: "dnf",
		Packages: []string{
			"python3-pip",
			"util-linux-user",
			"wget",
			"neovim",
			"zsh",
			"docker",
			"gh",
		},
	}
}

func (conf *PackageInstallerConfig) Fill(filePath string) {
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
