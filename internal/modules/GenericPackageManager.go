package modules

import (
	"errors"
	"fmt"
	"github.com/Linkinlog/gasible/internal"
	"github.com/Linkinlog/gasible/internal/app"
	"gopkg.in/yaml.v3"
	"log"
)

type config struct {
	Enabled        bool                   `yaml:"enabled"`
	ConfigSettings PackageManagerSettings `yaml:"settings"`
}

type GenericPackageManager struct {
	Name              string
	config            config
	PackageManagerMap map[packageManager][]string
	Application       *app.App
}

type PackageManagerSettings struct {
	Manager  string   `yaml:"manager"`
	Packages []string `yaml:"packages"`
}

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

func (gpm *GenericPackageManager) Setup() error {
	gpm.PackageManagerMap[gpm.Manager()] = gpm.config.ConfigSettings.Packages
	return gpm.managePackages("install")
}

func (gpm *GenericPackageManager) TearDown() error {
	gpm.PackageManagerMap[gpm.Manager()] = gpm.config.ConfigSettings.Packages
	return gpm.managePackages("uninstall")
}

func (gpm *GenericPackageManager) Update() error {
	gpm.PackageManagerMap[gpm.Manager()] = gpm.config.ConfigSettings.Packages
	return gpm.managePackages("update")
}

func (gpm *GenericPackageManager) GetName() string {
	return gpm.Name
}

func (gpm *GenericPackageManager) Config() app.ModuleConfig {
	return app.ModuleConfig{
		Enabled:  gpm.config.Enabled,
		Settings: gpm.config.ConfigSettings,
	}
}

func (gpm *GenericPackageManager) SetApp(app *app.App) {
	gpm.Application = app
}

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

type packageManager interface {
	Install([]string, *SysCall) error
	Uninstall([]string, *SysCall) error
	Update([]string, *SysCall) error
}

type BasePackageManager struct {
	Name string
	Args packageManagerArgs
	Opts packageManagerOpts
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

func (pm *BasePackageManager) Install(packages []string, call *SysCall) error {
	return pm.execute("install", packages, *call)
}

func (pm *BasePackageManager) Uninstall(packages []string, call *SysCall) error {
	return pm.execute("uninstall", packages, *call)
}

func (pm *BasePackageManager) Update(packages []string, call *SysCall) error {
	return pm.execute("update", packages, *call)
}

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

func (gpm *GenericPackageManager) Manager() *BasePackageManager {
	return supportedPackageManagers[gpm.config.ConfigSettings.Manager]
}

func (gpm *GenericPackageManager) addToInstaller(pm packageManager, packages []string) {
	if _, ok := gpm.PackageManagerMap[pm]; !ok {
		gpm.PackageManagerMap[pm] = []string{}
	}
	gpm.PackageManagerMap[pm] = append(gpm.PackageManagerMap[pm], packages...)
}

func (gpm *GenericPackageManager) system() *SysCall {
	sysCallMod := gpm.Application.ModuleRegistry.GetModule("SysCall")
	return sysCallMod.(*SysCall)
}

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

// PackageManagerMap
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
