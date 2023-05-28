package modules

import (
	"errors"
	"github.com/Linkinlog/gasible/internal/core"
	"gopkg.in/yaml.v3"
	"log"
	"runtime"
)

// Variable declaration

type GenericPackageManager struct {
	Priority int
	Enabled  bool
	Settings PackageManagerConfig
}

type PackageManagerConfig struct {
	Manager  string   `yaml:"pkg-manager"`
	Packages []string `yaml:"packages"`
}

type PackageManager interface {
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
	core.ModuleRegistry.Register("GenericPackageManager", &GenericPackageManager{
		Priority: 0,
		Enabled:  true,
		Settings: PackageManagerConfig{},
	})
}

// Interface methods

func (packageMan *GenericPackageManager) Setup() error {
	return InstallPackages(packageMan.Settings.Packages)
}

func (packageMan *GenericPackageManager) TearDown() error {
	return UninstallPackages(packageMan.Settings.Packages)
}

func (packageMan *GenericPackageManager) Update() error {
	return UpdatePackages(packageMan.Settings.Packages)
}

func (packageMan *GenericPackageManager) Config() core.ModuleConfig {
	return core.ModuleConfig{
		Priority: packageMan.Priority,
		Enabled:  packageMan.Enabled,
		Settings: packageMan.Settings,
	}
}

func (packageMan *GenericPackageManager) ParseSettings(rawSettings map[string]interface{}) (err error) {
	settingsBytes, err := yaml.Marshal(rawSettings)
	if err != nil {
		return
	}

	var settings PackageManagerConfig
	err = yaml.Unmarshal(settingsBytes, &settings)
	if err != nil {
		return
	}
	packageMan.Settings = settings
	return nil
}

// Methods that may be useful for other packages

func UpdatePackages(packages []string) (err error) {
	return managePackages(packages, "update")
}

func InstallPackages(packages []string) (err error) {
	return managePackages(packages, "install")
}

func UninstallPackages(packages []string) (err error) {
	return managePackages(packages, "uninstall")
}

// Helper functions

func managePackages(packages []string, operation string) (err error) {
	sys := core.System{
		Name:   runtime.GOOS,
		Runner: core.SudoRunner{},
	}

	moduleSettings, err := core.ModuleRegistry.Get("GenericPackageManager")
	if err != nil || moduleSettings == nil {
		log.Fatal("Failed to get GenericPackageManager module settings.")
	}

	packageMgr := determinePackageMgr(sys.Name, moduleSettings.Config().Settings.(PackageManagerConfig).Manager)

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

	out, err := sys.Exec(packageMgr.getExecutable(), packagesAndArgs)
	if err != nil {
		return errors.Join(err, errors.New(string(out)))
	}
	log.Printf("Package %s finished.\n Output: %s\n", operation, string(out))
	return
}

func determinePackageMgr(os string, manager string) (packageMgr PackageManager) {
	if os == "darwin" {
		// Failsafe as we only support brew on Mac
		// TODO maybe?
		return packageManagerMap["brew"]
	}
	return packageManagerMap[manager]
}

// formatCommand
// Should format the shell command
// with the proper operation (install, update, etc).
func formatCommand(packageMgr PackageManager, operation string) []string {
	return []string{
		operation,
		packageMgr.getCommandOptions().AutoConfirmOpt,
		packageMgr.getCommandOptions().QuietOpt,
	}
}

// packageManagerMap
// Give it a string, get a PackageManager.
var packageManagerMap = map[string]PackageManager{
	"apt":      &aptitude{},
	"apt-get":  &aptitude{},
	"aptitude": &aptitude{},
	"brew":     &brew{},
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
		UpdateArg:    "update",
		UpgradeArg:   "upgrade",
	}
}

func (apt *aptitude) getCommandOptions() *packageManagerOpts {
	return &packageManagerOpts{
		AutoConfirmOpt: "-y",
		QuietOpt:       "-qq",
	}
}
