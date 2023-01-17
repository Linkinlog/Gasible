package models

import (
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

// PackageInstallerConfig holds all fields fields
// relative to the package installer service.
type PackageInstallerConfig struct {
	Manager  string   `yaml:"pkg-manager-command,omitempty"`
	Args     string   `yaml:"command-args,omitempty"`
	Packages []string `yaml:"packages"`
}

// Map to validate if a package manager is supported.
var supportedPM = map[string]bool{
	"apt-get": true,
	"yum":     true,
	"dnf":     true,
	"zypper":  true,
	"pacman":  true,
	"winget":  true,
}

// Check if the package manager is supported,
// and if so, return the full path to it.
func (p PackageInstallerConfig) CheckPMAndReturnPath() string {
	pm := p.Manager
	if _, ok := supportedPM[pm]; !ok {
		err := fmt.Sprintf("Error: Package manager %s not found.", pm)
		panic(err)
	}
	path, err := exec.LookPath(pm)
	if err != nil {
		err := fmt.Sprintf("Error: Package manager %s not found.", pm)
		panic(err)
	}
	if os.Geteuid() != 0 {
		// TODO handle this better
		//panic("Error: Permission denied.")
	}
	return path
}

// Populate the struct with the default config for the package installer.
func (PackageInstallerConfig) Default() *PackageInstallerConfig {
	return &PackageInstallerConfig{
		Manager: "dnf",
		Args:    "install -y",
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

