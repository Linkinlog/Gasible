package core

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// PackageInstallerConfig holds all fields fields
// relative to the package installer service.
type PackageInstallerConfig struct {
	Manager  string   `yaml:"pkg-manager-command,omitempty"`
	Packages []string `yaml:"packages"`
}

type PackageInstallerOpts struct {
	AutoConfirmOpt string
	QuietOpt       string
}

type PackageInstallerSubCmds struct {
	InstallSubCmd string
	RemoveSubCmd  string
	UpdateSubCmd  string
	UpgradeSubCmd string
}

type PackageInstaller interface {
	GetExecutablePath(installerConfig *PackageInstallerConfig) string
	GetPackages(installerConfig *PackageInstallerConfig) []string
	GetSubCommands() *PackageInstallerSubCmds
	GetCommandOptions() *PackageInstallerOpts
}

// var Aptitude = PackageInstaller{
// 	AutoConfirmOpt: "-yy",
// 	QuietOpt:       "-q",
// 	InstallCmd:     "install",
// 	RemoveCmd:      "remove",
// 	UpdateCmd:      "update",
// 	UpgradeCmd:     "upgrade",
// }

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
func CheckPMAndReturnPath(pkgManager string) (string, error) {
	pm := pkgManager
	if _, ok := supportedPM[pm]; !ok {
		err := fmt.Sprintf("Error: Package manager %s not supported.", pm)
		return "", errors.New(err)
	}
	path, err := exec.LookPath(pm)
	if err != nil {
		err := fmt.Sprintf("Error: Package manager %s not found.", pm)
		return "", errors.New(err)
	}
	//if os.Geteuid() != 0 {
	// TODO handle this better
	//panic("Error: Permission denied.")
	//}
	return path, nil
}

// GetCmd returns a formatted string to install pkgs.
// It contains the pkg managers full path and arguments,
// and all the packages for it to install.
func (c PackageInstallerConfig) GetCmd() (string, error) {
	// Validate the package manager and get its root path.
	pm, err := CheckPMAndReturnPath(c.Manager)
	if err != nil {
		return "", err
	}
	// Turn our slice of packages into a single string.
	packages := strings.Join(c.Packages, " ")
	// Format all of the above into a string.
	command := strings.Join([]string{pm, c.Args, packages}, " ")
	return command, nil
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
