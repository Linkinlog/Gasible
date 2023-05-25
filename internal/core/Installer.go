package core

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// PackageManagerConfig holds all fields
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

// Opts InstallerOpts is the relative options that we
// will need to pass into Run.
type InstallerOpts struct {
	NoOp bool
	Os   *System
}

// Run installs the packages listed in the
// packages section of the YAML file.
//
// Deprecated: Use GenericPackageManager.
func (opts *InstallerOpts) Run(c *PackageManagerConfig) ([]byte, error) {
	// Format our command
	command, err := c.GetInstallCommand()
	if err != nil {
		return []byte{}, err
	}

	return opts.Os.Exec(command)
}

// GetManagerPath Check if the package manager is in $PATH,
// and if so, return the full path to it.
func (pkgManagerConf *PackageManagerConfig) GetManagerPath() (string, error) {
	path, err := exec.LookPath(pkgManagerConf.Manager)
	if err != nil {
		err := fmt.Sprintf(
			"Error: Package manager %s not found.",
			pkgManagerConf.Manager,
		)
		return "", errors.New(err)
	}
	//if os.Geteuid() != 0 {
	// TODO handle this better
	//panic("Error: Permission denied.")
	//}
	return path, nil
}

func InstallWithDetectedOS(packages []string) error {
	conf := GetConfig()
	installerOpts := InstallerOpts{
		NoOp: conf.GlobalOpts.NoOp,
		Os: &System{
			Name:   runtime.GOOS,
			Runner: RealRunner{},
		},
	}
	_, err := installerOpts.Run(&conf.PackageManagerConfig)
	if err != nil {
		return err
	}
	return nil
}

// GetInstallCommand returns a formatted string to install pkgs.
// It contains the pkg managers full path and arguments,
// as well as all the packages for it to install.
func (pkgManagerConf *PackageManagerConfig) GetInstallCommand() (string, error) {
	// Validate the package manager and get its root path.
	pkgManager, err := pkgManagerConf.GetManagerPath()
	if err != nil {
		return "", err
	}
	// Format all of the above into a string.
	command := strings.Join([]string{
		pkgManager,
		pkgManagerConf.ChosenManager.GetSubCommands().InstallArg,
		strings.Join(pkgManagerConf.Packages, " "),
		pkgManagerConf.ChosenManager.GetCommandOptions().AutoConfirmOpt,
		pkgManagerConf.ChosenManager.GetCommandOptions().QuietOpt,
	}, " ")
	return command, nil
}

// Default Populate the struct with the default config for the package installer.
func (*PackageManagerConfig) Default() *PackageManagerConfig {
	return &PackageManagerConfig{
		Manager: "apt",
		Packages: []string{
			"neovim",
		},
	}
}
