package yamlparser

import (
	"os"

	"gopkg.in/yaml.v3"
)

type PackageInstallerConfig struct {
	Manager  string   `yaml:"pkg-manager,omitempty"`
	Packages []string `yaml:"packages,omitempty"`
}

type ServicesConfig struct {
	Teamviewer bool `yaml:"teamviewer,omitempty"`
	Ssh        bool `yaml:"ssh,omitempty"`
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
	PackageInstallerConfig `yaml:",inline"`
	ServicesConfig         `yaml:",inline"`
	GeneralConfig          `yaml:",inline"`
}

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
		Teamviewer: true,
		Ssh:        true,
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
