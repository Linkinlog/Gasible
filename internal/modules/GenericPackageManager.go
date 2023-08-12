package modules

import (
	"errors"
	"fmt"
	"github.com/Linkinlog/gasible/internal"
	"github.com/Linkinlog/gasible/internal/app"
	"gopkg.in/yaml.v3"
	"log"
)

// init
// This should really just handle registering the module in the registry.
func init() {
	ToBeRegistered = append(ToBeRegistered, &GenericPackageManager{
		Name: "GenericPackageManager",
		config: config{
			Enabled:        true,
			ConfigSettings: PackageManagerSettings{},
		},
		PackageManagerMap: make(map[packageManager][]string),
	})
}

// GenericPackageManager implements module, so we can register and execute it.
type GenericPackageManager struct {
	Name              string
	config            config
	PackageManagerMap map[packageManager][]string
	Application       *app.App
}

// config is the YAML configuration for GenericPackageManager.
type config struct {
	Enabled        bool                   `yaml:"enabled"`
	ConfigSettings PackageManagerSettings `yaml:"settings"`
}

// PackageManagerSettings contains the user chosen package manager and the packages the user wants to install.
type PackageManagerSettings struct {
	Manager  string   `yaml:"manager"`
	Packages []string `yaml:"packages"`
}

// Setup will run the installation command on the chosen package manager.
func (gpm *GenericPackageManager) Setup() error {
	gpm.PackageManagerMap[gpm.Manager()] = gpm.config.ConfigSettings.Packages
	return gpm.managePackages("install")
}

// TearDown will run the remove command on the chosen package manager.
func (gpm *GenericPackageManager) TearDown() error {
	gpm.PackageManagerMap[gpm.Manager()] = gpm.config.ConfigSettings.Packages
	return gpm.managePackages("uninstall")
}

// Update will run the update command on the chosen package manager.
func (gpm *GenericPackageManager) Update() error {
	gpm.PackageManagerMap[gpm.Manager()] = gpm.config.ConfigSettings.Packages
	return gpm.managePackages("update")
}

// GetName returns the name field of the GenericPackageManager struct.
func (gpm *GenericPackageManager) GetName() string {
	return gpm.Name
}

// Config returns the shallow-copied module config from our module's config.
func (gpm *GenericPackageManager) Config() app.ModuleConfig {
	return app.ModuleConfig{
		Enabled:  gpm.config.Enabled,
		Settings: gpm.config.ConfigSettings,
	}
}

// SetApp sets the application field as the app that is passed in.
func (gpm *GenericPackageManager) SetApp(app *app.App) {
	gpm.Application = app
}

// ParseConfig takes in a map that ideally contains a YAML structure, to be marshalled into the config.
func (gpm *GenericPackageManager) ParseConfig(rawConfig map[string]interface{}) error {
	configBytes, err := yaml.Marshal(rawConfig)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configBytes, &gpm.config)
	if err != nil {
		return err
	}
	return nil
}

// packageManager is an interface that is meant to group the functions that a package manager would need to do.
type packageManager interface {
	Install([]string, *SysCall) error
	Uninstall([]string, *SysCall) error
	Update([]string, *SysCall) error
}

// BasePackageManager implements the packageManager interface.
type BasePackageManager struct {
	Name string
	Args packageManagerArgs
	Opts packageManagerOpts
}

// packageManagerArgs contains what we need to tell each supported package manager what we intend to do.
type packageManagerArgs struct {
	InstallArg   string
	UninstallArg string
	UpdateArg    string
	UpgradeArg   string
}

// packageManagerOpts contains what we need to skip prompts and quiet the output when there are no errors.
type packageManagerOpts struct {
	AutoConfirmOpt string
	QuietOpt       string
}

// Install install the packages.
func (pm *BasePackageManager) Install(packages []string, call *SysCall) error {
	return pm.execute("install", packages, *call)
}

// Uninstall uninstall the packages.
func (pm *BasePackageManager) Uninstall(packages []string, call *SysCall) error {
	return pm.execute("uninstall", packages, *call)
}

// Update updates the packages.
func (pm *BasePackageManager) Update(packages []string, call *SysCall) error {
	return pm.execute("update", packages, *call)
}

// execute ensures the package manager is set, checks if we need sudo, formats the command, and manages the packages.
func (pm *BasePackageManager) execute(operation string, packages []string, syscall SysCall) error {
	var noPackageManagerFoundErr = errors.New("no package Manager set, set one in the config")
	if pm == nil {
		return internal.ErrorAs("basePackageManager.execute", noPackageManagerFoundErr)
	}
	if len(packages) < 1 {
		return nil
	}
	sudo := pm.Name != "brew"
	formattedCommand := formatCommand(pm, operation)
	packagesAndArgs := append(formattedCommand, packages...)
	out, execErr := syscall.Exec(pm.Name, packagesAndArgs, sudo)
	if execErr != nil {
		return fmt.Errorf("%w: %s", execErr, string(out))
	}
	log.Printf("%d Packages finished running operation: %s\n", len(packages), operation)
	return nil
}

// formatCommand will set the proper auto-confirm and quiet options on our management command.
func formatCommand(pm *BasePackageManager, operation string) []string {
	var args string
	switch operation {
	case "install":
		args = pm.Args.InstallArg
	case "uninstall":
		args = pm.Args.UninstallArg
	case "update":
		args = pm.Args.UpdateArg
	case "upgrade":
		args = pm.Args.UpgradeArg
	}

	return []string{
		args,
		pm.Opts.AutoConfirmOpt,
		pm.Opts.QuietOpt,
	}
}

// Manager will get the current package manager as long as it is supprorted.
func (gpm *GenericPackageManager) Manager() *BasePackageManager {
	return supportedPackageManagers[gpm.config.ConfigSettings.Manager]
}

// system returns the SysCall in use by the registry.
func (gpm *GenericPackageManager) system() *SysCall {
	sysCallMod := gpm.Application.ModuleRegistry.GetModule("SysCall")
	return sysCallMod.(*SysCall)
}

// managePackages will take an operation such as "install" and install all packages in the PackageManagerMap.
func (gpm *GenericPackageManager) managePackages(operation string) error {
	for pm, packages := range gpm.PackageManagerMap {
		switch operation {
		case "install":
			if err := pm.Install(packages, gpm.system()); err != nil {
				return err
			}
		case "uninstall":
			if err := pm.Uninstall(packages, gpm.system()); err != nil {
				return err
			}
		case "update":
			if err := pm.Update(packages, gpm.system()); err != nil {
				return err
			}
		}
	}
	return nil
}

// supportedPackageManagers
// Give it a string, get a BasePackageManager.
var supportedPackageManagers = map[string]*BasePackageManager{
	"apt":      &aptitude,
	"apt-get":  &aptitude,
	"Aptitude": &aptitude,
	"brew":     &brew,
	"dnf":      &dnf,
	"pacman":   &pacman,
	"zypper":   &zypper,
}

// package Manager methods &structs are below

//func (pm *basePackageManager) getExecutable() string {
//	return pm.Name
//}

// brew is the package Manager for Mac
var brew = BasePackageManager{
	Name: "brew",
	Args: packageManagerArgs{
		InstallArg:   "install",
		UninstallArg: "uninstall",
		UpdateArg:    "update",
		UpgradeArg:   "upgrade",
	},
	Opts: packageManagerOpts{
		AutoConfirmOpt: "",
		QuietOpt:       "-q",
	},
}

// aptitude // apt-get // apt is for debian based distros
var aptitude = BasePackageManager{
	Name: "apt-get",
	Args: packageManagerArgs{
		InstallArg:   "install",
		UninstallArg: "remove",
		UpdateArg:    "install",
		UpgradeArg:   "upgrade",
	},
	Opts: packageManagerOpts{
		AutoConfirmOpt: "-y",
		QuietOpt:       "-qq",
	},
}

// dnf is for RPM / Redhat-like distros
var dnf = BasePackageManager{
	Name: "dnf",
	Args: packageManagerArgs{
		InstallArg:   "install",
		UninstallArg: "remove",
		UpdateArg:    "upgrade",
		UpgradeArg:   "upgrade",
	},
	Opts: packageManagerOpts{
		AutoConfirmOpt: "-y",
		QuietOpt:       "-q",
	},
}

// pacman is for arch
var pacman = BasePackageManager{
	Name: "pacman",
	Args: packageManagerArgs{
		InstallArg:   "-S",
		UninstallArg: "-R",
		UpdateArg:    "-Syu",
		UpgradeArg:   "-Syu",
	},
	Opts: packageManagerOpts{
		AutoConfirmOpt: "--noconfirm",
		QuietOpt:       "--quiet",
	},
}

// zypper is for Suse
var zypper = BasePackageManager{
	Name: "zypper",
	Args: packageManagerArgs{
		InstallArg:   "in",
		UninstallArg: "rm",
		UpdateArg:    "in",
		UpgradeArg:   "up",
	},
	Opts: packageManagerOpts{
		AutoConfirmOpt: "--non-interactive",
		QuietOpt:       "--quiet",
	},
}
