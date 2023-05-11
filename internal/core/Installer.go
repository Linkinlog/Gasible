package core

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// PackageManagerConfig holds all fields fields
// relative to the package installer service.
type PackageManagerConfig struct {
	Manager       string         `yaml:"pkg-manager-command,omitempty"`
	Packages      []string       `yaml:"packages"`
	ChosenManager PackageManager `yaml:"-"`
}

type PackageManagerOpts struct {
	AutoConfirmOpt string
	QuietOpt       string
}

type PackageManagerArgs struct {
	InstallArg string
	RemoveArg  string
	UpdateArg  string
	UpgradeArg string
}

type PackageManager interface {
	GetSubCommands() *PackageManagerArgs
	GetCommandOptions() *PackageManagerOpts
}

// Check if the package manager is in $PATH,
// and if so, return the full path to it.
func (pkgInstallConf *PackageManagerConfig) GetManagerPath() (string, error) {
	path, err := exec.LookPath(pkgInstallConf.Manager)
	if err != nil {
		err := fmt.Sprintf(
			"Error: Package manager %s not found.",
			pkgInstallConf.Manager,
		)
		return "", errors.New(err)
	}
	//if os.Geteuid() != 0 {
	// TODO handle this better
	//panic("Error: Permission denied.")
	//}
	return path, nil
}

// GetInstallCommand returns a formatted string to install pkgs.
// It contains the pkg managers full path and arguments,
// as well as all the packages for it to install.
func (c PackageManagerConfig) GetInstallCommand() (string, error) {
	// Validate the package manager and get its root path.
	pkgManager, err := c.GetManagerPath()
	if err != nil {
		return "", err
	}
	// Format all of the above into a string.
	command := strings.Join([]string{
		pkgManager,
		c.ChosenManager.GetSubCommands().InstallArg,
		strings.Join(c.Packages, " "),
		c.ChosenManager.GetCommandOptions().AutoConfirmOpt,
		c.ChosenManager.GetCommandOptions().QuietOpt,
	}, " ")
	return command, nil
}

// Populate the struct with the default config for the package installer.
func (PackageManagerConfig) Default() *PackageManagerConfig {
	return &PackageManagerConfig{
		Manager: "apt",
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
