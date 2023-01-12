package yamlParser

import (
	"os"

	"gopkg.in/yaml.v3"
)

type PackageInstallerConfig struct {
	Manager  string   `yaml:"pkg-manager,omitempty"`
	Packages []string `yaml:"packages,omitempty"`
}

type ServicesConfig struct {
	Installer  bool `yaml:"installer,omitempty"`
	Teamviewer bool `yaml:"teamviewer,omitempty"`
	Ssh        bool `yaml:"ssh,omitempty"`
	Git        bool `yaml:"git,omitempty"`
}

type GeneralConfig struct {
	Hostname        string  `yaml:"hostname,omitempty"`
	IP              string  `yaml:"staticIP,omitempty"`
	Mask            string  `yaml:"mask,omitempty"`
	TeamViewerCreds TVCreds `yaml:"TeamViewerCreds,omitempty"`
}

type TVCreds struct {
	User string `yaml:"user,omitempty"`
	Pass string `yaml:"pass,omitempty"`
}

type Config struct {
	PackageInstallerConfig `yaml:",inline,omitempty"`
	ServicesConfig         `yaml:",inline,omitempty"`
	GeneralConfig          `yaml:",inline,omitempty"`
}

type IConfig interface {
    Parse() error
    Defaults() error
}

// CreateDefaults will generate a YAML file
// using the defaults we outline.
func CreateDefaults() error {
	pkgInstallConf := PackageInstallerConfig{
		Manager: "dnf",
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
	Services := ServicesConfig{
        Installer:  true,
		Teamviewer: true,
		Ssh:        true,
		Git:        true,
	}
	General := GeneralConfig{
		Hostname: "development-station",
		IP:       "192.168.4.20",
		Mask:     "255.255.255.0",
		TeamViewerCreds: TVCreds{
			User: "username",
			Pass: "password",
		},
	}
	Conf := Config{
		pkgInstallConf,
		Services,
		General,
	}
	d, err := yaml.Marshal(&Conf)
	if err != nil {
		return err
	}

	err = os.WriteFile("default.yml", d, 0644)
	if err != nil {
		return err
	}

	return nil
}
