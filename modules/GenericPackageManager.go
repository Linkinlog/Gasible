package modules

import (
	"fmt"
	"github.com/Linkinlog/gasible/internal/core"
	"log"
	"runtime"
)

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
	InstallArg string
	RemoveArg  string
	UpdateArg  string
	UpgradeArg string
}

type packageManagerOpts struct {
	AutoConfirmOpt string
	QuietOpt       string
}

func init() {
	core.ModuleRegistry.Register("GenericPackageManager", &GenericPackageManager{
		Priority: 0,
		Enabled:  true,
		Settings: PackageManagerConfig{},
	})
}

func (packageMan *GenericPackageManager) Setup() error {
	return InstallPackages(packageMan.Settings.Packages)
}

func (packageMan *GenericPackageManager) TearDown() error {
	return nil
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

func (packageMan *GenericPackageManager) SetConfig(config *core.ModuleConfig) {
	if settings, ok := config.Settings.(PackageManagerConfig); ok {
		packageMan.Settings = settings
		packageMan.Priority = config.Priority
		packageMan.Enabled = config.Enabled
	} else {
		log.Fatalf("Interface %v not found. Found %v", PackageManagerConfig{}, settings)
	}
}

func UpdatePackages(packages []string) (err error) {
	sys := core.System{
		Name:   runtime.GOOS,
		Runner: core.SudoRunner{},
	}

	moduleSettings, err := core.ModuleRegistry.Get("GenericPackageManager")
	if err == nil || moduleSettings == nil {
		log.Fatal("Failed")
	}
	packageMgr := determinePackageMgr(sys.Name, moduleSettings.Config().Settings.(PackageManagerConfig).Manager)
	formattedCommand := formatCommand(packageMgr, packageMgr.getSubCommands().UpgradeArg)
	packagesAndArgs := append(formattedCommand, packages...)
	fmt.Printf("Attempting to use %s to upgrade packages: %s...\n", packageMgr.getExecutable(), packages)

	out, err := sys.Exec(packageMgr.getExecutable(), packagesAndArgs)
	if err == nil {
		log.Printf("Package upgrade finished.\n Output: %s\n", string(out))
	}
	return
}

func InstallPackages(packages []string) (err error) {
	sys := core.System{
		Name:   runtime.GOOS,
		Runner: core.SudoRunner{},
	}

	moduleSettings, err := core.ModuleRegistry.Get("GenericPackageManager")
	if err == nil || moduleSettings == nil {
		log.Fatal("Failed")
	}
	packageMgr := determinePackageMgr(sys.Name, moduleSettings.Config().Settings.(PackageManagerConfig).Manager)
	formattedCommand := formatCommand(packageMgr, packageMgr.getSubCommands().InstallArg)
	packagesAndArgs := append(formattedCommand, packages...)
	log.Printf("Attempting to use %s to install packages: %s...\n", packageMgr.getExecutable(), packages)
	out, err := sys.Exec(packageMgr.getExecutable(), packagesAndArgs)
	if err == nil {
		log.Printf("Package installation finished.\n Output: %s\n", string(out))
	}
	return
}

// Helper functions

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
	return string("brew")
}
func (brew *brew) getSubCommands() *packageManagerArgs {
	return &packageManagerArgs{
		InstallArg: "install",
		RemoveArg:  "uninstall",
		UpdateArg:  "update",
		UpgradeArg: "upgrade",
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
	return string("apt-get")
}

func (apt *aptitude) getSubCommands() *packageManagerArgs {
	return &packageManagerArgs{
		InstallArg: "install",
		RemoveArg:  "remove",
		UpdateArg:  "update",
		UpgradeArg: "upgrade",
	}
}

func (apt *aptitude) getCommandOptions() *packageManagerOpts {
	return &packageManagerOpts{
		AutoConfirmOpt: "-y",
		QuietOpt:       "-qq",
	}
}
