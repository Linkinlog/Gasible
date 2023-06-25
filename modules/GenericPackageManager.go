package modules

import (
	"errors"
	"fmt"
	"github.com/Linkinlog/gasible/internal/core"
	"gopkg.in/yaml.v3"
	"log"
)

// Variable declaration

type genericPackageManager struct {
	name     string
	Enabled  bool
	Settings packageManagerConfig
}

type packageManagerConfig struct {
	Manager      string   `yaml:"manager"`
	Packages     []string `yaml:"packages"`
	Dependencies []string `yaml:"dependencies"`
}

type packageManager interface {
	getExecutable() string
	getSubCommands() *packageManagerArgs
	getCommandOptions() *packageManagerOpts
}

type packageManagerArgs struct {
	InstallArg   string
	UninstallArg string
	UpdateArg    string
	UpgradeArg   string
}

type packageManagerOpts struct {
	AutoConfirmOpt string
	QuietOpt       string
}

// init
// This should really just handle registering the module in the registry.
func init() {
	core.ModuleRegistry.Register(&genericPackageManager{
		name:     "GenericPackageManager",
		Enabled:  true,
		Settings: packageManagerConfig{},
	})
}

// Interface methods

func (packageMan *genericPackageManager) Setup() error {
	return managePackages(packageMan.Settings.Packages, "install")
}

func (packageMan *genericPackageManager) TearDown() error {
	return managePackages(packageMan.Settings.Packages, "uninstall")
}

func (packageMan *genericPackageManager) Update() error {
	return managePackages(packageMan.Settings.Packages, "update")
}

func (packageMan *genericPackageManager) Name() string {
	return packageMan.name
}

func (packageMan *genericPackageManager) Config() core.ModuleConfig {
	return core.ModuleConfig{
		Enabled:  packageMan.Enabled,
		Settings: packageMan.Settings,
	}
}

func (packageMan *genericPackageManager) ParseConfig(rawConfig map[string]interface{}) (err error) {
	configBytes, err := yaml.Marshal(rawConfig)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(configBytes, packageMan)
	if err != nil {
		return
	}
	return nil
}

// Methods that may be useful for other packages

// addToInstaller is a helper function so other modules can install packages.
// It takes packageMap which is a map of a package manager name to a slice of packages to install.
// Allows someone to create a map of all supported package managers and the differing packages between them.
func (packageMan *genericPackageManager) addToInstaller(packageMap map[string][]string) error {
	moduleSettings, err := core.ModuleRegistry.Get("GenericPackageManager")
	if err != nil || moduleSettings == nil {
		return errors.New("failed to get GenericPackageManager module settings")
	}
	packageMgr, err := determinePackageMgr(packageMan.Settings.Manager)
	if err != nil {
		return err
	}
	if deps, ok := packageMap[packageMgr.getExecutable()]; !ok {
		return fmt.Errorf("failed to find package manager %s in package map", packageMgr)
	} else {
		packageMan.Settings.Packages = append(packageMan.Settings.Packages, deps...)
		return nil
	}
}

// Helper functions

func managePackages(packages []string, operation string) (err error) {
	moduleSettings, err := core.ModuleRegistry.Get("GenericPackageManager")
	if err != nil || moduleSettings == nil {
		return errors.New("failed to get GenericPackageManager module settings")
	}

	packageMgr, err := determinePackageMgr(moduleSettings.Config().Settings.(packageManagerConfig).Manager)
	if err != nil {
		return err
	}

	// Get the appropriate command argument based on the operation type.
	var commandArg string
	switch operation {
	case "install":
		commandArg = packageMgr.getSubCommands().InstallArg
	case "uninstall":
		commandArg = packageMgr.getSubCommands().UninstallArg
	case "update":
		commandArg = packageMgr.getSubCommands().UpdateArg
	}

	formattedCommand := formatCommand(packageMgr, commandArg)
	packagesAndArgs := append(formattedCommand, packages...)
	log.Printf("Attempting to use %s to %s packages: %s...\n", packageMgr.getExecutable(), operation, packages)

	out, err := core.CurrentState.System.Exec(packageMgr.getExecutable(), packagesAndArgs)
	if err != nil {
		return errors.Join(err, errors.New(string(out)))
	}
	log.Printf("%sPackage %s finished.\n", string(out), operation)
	return
}

func determinePackageMgr(manager string) (packageMgr packageManager, err error) {
	var ok bool
	os := core.CurrentState.System.Name
	if os == "darwin" {
		// Failsafe as we only support brew on Mac.
		// Also, brew doesn't support being ran as sudo.
		// TODO maybe?
		packageMgr, ok = packageManagerMap["brew"]
	} else {
		packageMgr, ok = packageManagerMap[manager]
		core.CurrentState.System.Runner = core.SudoRunner{}
	}
	if ok {
		return packageMgr, nil
	} else {
		return nil, fmt.Errorf("package manager %s not found", manager)
	}
}

// formatCommand
// Should format the shell command
// with the proper operation (install, update, etc).
func formatCommand(packageMgr packageManager, operation string) []string {
	return []string{
		operation,
		packageMgr.getCommandOptions().AutoConfirmOpt,
		packageMgr.getCommandOptions().QuietOpt,
	}
}

// packageManagerMap
// Give it a string, get a packageManager.
var packageManagerMap = map[string]packageManager{
	"apt":      &aptitude{},
	"apt-get":  &aptitude{},
	"aptitude": &aptitude{},
	"brew":     &brew{},
	"dnf":      &dnf{},
	"pacman":   &pacman{},
	"zypper":   &zypper{},
}

// Package manager structs below

// brew
type brew struct{}

func (brew *brew) getExecutable() string {
	return "brew"
}
func (brew *brew) getSubCommands() *packageManagerArgs {
	return &packageManagerArgs{
		InstallArg:   "install",
		UninstallArg: "uninstall",
		UpdateArg:    "update",
		UpgradeArg:   "upgrade",
	}
}

func (brew *brew) getCommandOptions() *packageManagerOpts {
	return &packageManagerOpts{
		AutoConfirmOpt: "",
		QuietOpt:       "-q",
	}
}

// aptitude // apt-get // apt
type aptitude struct{}

func (apt *aptitude) getExecutable() string {
	return "apt-get"
}

func (apt *aptitude) getSubCommands() *packageManagerArgs {
	return &packageManagerArgs{
		InstallArg:   "install",
		UninstallArg: "remove",
		UpdateArg:    "install",
		UpgradeArg:   "upgrade",
	}
}

func (apt *aptitude) getCommandOptions() *packageManagerOpts {
	return &packageManagerOpts{
		AutoConfirmOpt: "-y",
		QuietOpt:       "-qq",
	}
}

// dnf
type dnf struct{}

func (dnf *dnf) getExecutable() string {
	return "dnf"
}

func (dnf *dnf) getSubCommands() *packageManagerArgs {
	return &packageManagerArgs{
		InstallArg:   "install",
		UninstallArg: "remove",
		UpdateArg:    "upgrade",
		UpgradeArg:   "upgrade",
	}
}

func (dnf *dnf) getCommandOptions() *packageManagerOpts {
	return &packageManagerOpts{
		AutoConfirmOpt: "-y",
		QuietOpt:       "-q",
	}
}

// pacman
type pacman struct{}

func (pacman *pacman) getExecutable() string {
	return "pacman"
}

func (pacman *pacman) getSubCommands() *packageManagerArgs {
	return &packageManagerArgs{
		InstallArg:   "-S",
		UninstallArg: "-R",
		UpdateArg:    "-Syu",
		UpgradeArg:   "-Syu",
	}
}

func (pacman *pacman) getCommandOptions() *packageManagerOpts {
	return &packageManagerOpts{
		AutoConfirmOpt: "--noconfirm",
		QuietOpt:       "--quiet",
	}
}

// zypper
type zypper struct{}

func (zypper *zypper) getExecutable() string {
	return "zypper"
}

func (zypper *zypper) getSubCommands() *packageManagerArgs {
	return &packageManagerArgs{
		InstallArg:   "in",
		UninstallArg: "rm",
		UpdateArg:    "in",
		UpgradeArg:   "up",
	}
}

func (zypper *zypper) getCommandOptions() *packageManagerOpts {
	return &packageManagerOpts{
		AutoConfirmOpt: "--non-interactive",
		QuietOpt:       "--quiet",
	}
}
