package modules

import (
	"runtime"
	"strings"

	"github.com/Linkinlog/gasible/internal/core"
)

type GenericPackageManager struct {
	Name string
}

// PackageManagerConfig holds all fields
// relative to the package installer service.
type PackageManagerConfig struct {
	Manager       string         `yaml:"pkg-manager-command,omitempty"`
	Packages      []string       `yaml:"packages"`
	ChosenManager PackageManager `yaml:"-"`
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
	core.ModuleRegistry.Register("GenericPackageManager", &GenericPackageManager{})
}

func (packageMan *GenericPackageManager) Setup() error {
	return InstallPackages(core.CurrentConfig.Packages)
}

func (packageMan *GenericPackageManager) Update() error {
	return UpdatePackages(core.CurrentConfig.Packages)
}

func UpdatePackages(packages []string) (err error) {
	sys := core.System{
		Name:   runtime.GOOS,
		Runner: core.RealRunner{},
	}

	packageMgr := determinePackageMgr(sys.Name, core.CurrentConfig.Manager)
	formattedCommand := formatCommand(packageMgr, packageMgr.getSubCommands().UpdateArg)
	_, err = sys.Exec(formattedCommand, packages...)
	return
}

func InstallPackages(packages []string) (err error) {
	sys := core.System{
		Name:   runtime.GOOS,
		Runner: core.RealRunner{},
	}

	packageMgr := determinePackageMgr(sys.Name, core.CurrentConfig.Manager)
	formattedCommand := formatCommand(packageMgr, packageMgr.getSubCommands().InstallArg)
	_, err = sys.Exec(formattedCommand, packages...)
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
func formatCommand(packageMgr PackageManager, operation string) string {
	return strings.Join([]string{
		string(packageMgr.getExecutable()),
		operation,
		packageMgr.getCommandOptions().AutoConfirmOpt,
		packageMgr.getCommandOptions().QuietOpt,
	}, " ")
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
