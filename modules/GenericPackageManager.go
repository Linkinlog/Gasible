package modules

import (
	"errors"
	"fmt"
	"github.com/Linkinlog/gasible/internal/core"
	"gopkg.in/yaml.v3"
	"log"
)

// Variable declaration

type PackageManagerMap map[*packageManager][]string

type GenericPackageManager struct {
	Name              string
	Enabled           bool
	Settings          PackageManagerConfig
	PackageManagerMap PackageManagerMap
}

type PackageManagerConfig struct {
	Manager  string   `yaml:"manager"`
	Packages []string `yaml:"packages"`
}

type packageManagerer interface {
	getExecutable() string
	getSubCommands() *packageManagerArgs
	getCommandOptions() *packageManagerOpts
}

type packageManager struct {
	Name string
	packageManagerArgs
	packageManagerOpts
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
	core.ModuleRegistry.Register(&GenericPackageManager{
		Name:              "GenericPackageManager",
		Enabled:           true,
		Settings:          PackageManagerConfig{},
		PackageManagerMap: make(PackageManagerMap),
	})
}

// Interface methods

func (packageMan *GenericPackageManager) Setup() error {
	packageMan.validateAndAddModulePackages()
	return managePackages(packageMan.manager(), packageMan.Settings.Packages, "install")
}

func (packageMan *GenericPackageManager) TearDown() error {
	packageMan.validateAndAddModulePackages()
	return managePackages(packageMan.manager(), packageMan.Settings.Packages, "uninstall")
}

func (packageMan *GenericPackageManager) Update() error {
	packageMan.validateAndAddModulePackages()
	return managePackages(packageMan.manager(), packageMan.Settings.Packages, "update")
}

func (packageMan *GenericPackageManager) GetName() string {
	return packageMan.Name
}

func (packageMan *GenericPackageManager) Config() core.ModuleConfig {
	return core.ModuleConfig{
		Enabled:  packageMan.Enabled,
		Settings: packageMan.Settings,
	}
}

func (packageMan *GenericPackageManager) ParseConfig(rawConfig map[string]interface{}) (err error) {
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

// AddToInstaller is a helper function so other modules can install packages.
// It takes packageMap, which is a map of a package manager Name to a slice of packages to install.
// Allows someone to create a map of all supported package managers and the differing packages between them.
func (packageMan *GenericPackageManager) AddToInstaller(packageMap PackageManagerMap) {
	if len(packageMap) == 0 || packageMan.PackageManagerMap == nil {
		return
	}
	for key, val := range packageMap {
		packageMan.PackageManagerMap[key] = append(packageMan.PackageManagerMap[key], val...)
	}
}

// Helper functions

func (packageMan *GenericPackageManager) manager() *packageManager {
	moduleSettings, err := core.ModuleRegistry.Get("GenericPackageManager")
	if err != nil || moduleSettings == nil {
		log.Fatalf("failed to get GenericPackageManager module settings")
		return nil
	}

	packageMgr, err := determinePackageMgr(moduleSettings.Config().Settings.(PackageManagerConfig).Manager)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return packageMgr
}

func (packageMan *GenericPackageManager) validateAndAddModulePackages() {
	packagesForCurrentManager := packageMan.PackageManagerMap[packageMan.manager()]
	packageMan.Settings.Packages = append(packageMan.Settings.Packages, packagesForCurrentManager...)
}

func managePackages(packageMgr *packageManager, packages []string, operation string) (err error) {
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

func determinePackageMgr(manager string) (packageMgr *packageManager, err error) {
	var ok bool
	os := core.CurrentState.System.Name
	if os == "darwin" {
		// Failsafe as we only support brew on Mac.
		// Also, brew doesn't support being ran as sudo.
		// TODO maybe?
		packageMgr, ok = supportedPackageManagers["brew"]
	} else {
		packageMgr, ok = supportedPackageManagers[manager]
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
func formatCommand(packageMgr packageManagerer, operation string) []string {
	return []string{
		operation,
		packageMgr.getCommandOptions().AutoConfirmOpt,
		packageMgr.getCommandOptions().QuietOpt,
	}
}

// PackageManagerMap
// Give it a string, get a packageManager.
var supportedPackageManagers = map[string]*packageManager{
	"apt":      &Aptitude,
	"apt-get":  &Aptitude,
	"Aptitude": &Aptitude,
	"brew":     &Brew,
	"dnf":      &Dnf,
	"pacman":   &Pacman,
	"zypper":   &Zypper,
}

// package manager methods &structs are below

func (p *packageManager) getExecutable() string {
	return p.Name
}

func (p *packageManager) getSubCommands() *packageManagerArgs {
	return &p.packageManagerArgs
}

func (p *packageManager) getCommandOptions() *packageManagerOpts {
	return &p.packageManagerOpts
}

// Brew is the package manager for Mac
var Brew packageManager = packageManager{
	Name: "brew",
	packageManagerArgs: packageManagerArgs{
		InstallArg:   "install",
		UninstallArg: "uninstall",
		UpdateArg:    "update",
		UpgradeArg:   "upgrade",
	},
	packageManagerOpts: packageManagerOpts{
		AutoConfirmOpt: "",
		QuietOpt:       "-q",
	},
}

// Aptitude // apt-get // apt is for debian based distros
var Aptitude packageManager = packageManager{
	Name: "apt-get",
	packageManagerArgs: packageManagerArgs{
		InstallArg:   "install",
		UninstallArg: "remove",
		UpdateArg:    "install",
		UpgradeArg:   "upgrade",
	},
	packageManagerOpts: packageManagerOpts{
		AutoConfirmOpt: "-y",
		QuietOpt:       "-qq",
	},
}

// Dnf is for RPM / Redhat-like distros
var Dnf packageManager = packageManager{
	Name: "dnf",
	packageManagerArgs: packageManagerArgs{
		InstallArg:   "install",
		UninstallArg: "remove",
		UpdateArg:    "upgrade",
		UpgradeArg:   "upgrade",
	},
	packageManagerOpts: packageManagerOpts{
		AutoConfirmOpt: "-y",
		QuietOpt:       "-q",
	},
}

// Pacman is for arch
var Pacman packageManager = packageManager{
	Name: "pacman",
	packageManagerArgs: packageManagerArgs{
		InstallArg:   "-S",
		UninstallArg: "-R",
		UpdateArg:    "-Syu",
		UpgradeArg:   "-Syu",
	},
	packageManagerOpts: packageManagerOpts{
		AutoConfirmOpt: "--noconfirm",
		QuietOpt:       "--quiet",
	},
}

// Zypper is for Suse
var Zypper packageManager = packageManager{
	Name: "zypper",
	packageManagerArgs: packageManagerArgs{
		InstallArg:   "in",
		UninstallArg: "rm",
		UpdateArg:    "in",
		UpgradeArg:   "up",
	},
	packageManagerOpts: packageManagerOpts{
		AutoConfirmOpt: "--non-interactive",
		QuietOpt:       "--quiet",
	},
}
