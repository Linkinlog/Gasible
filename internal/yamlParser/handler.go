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
	PackageInstallerConfig `yaml:",inline"`
	ServicesConfig         `yaml:",inline"`
	GeneralConfig          `yaml:",inline"`
}

type IConfig interface {
    generate() error
    parse() error
}

type IPackageInstallerConfig interface {
    installPkgs() error
}

type IGeneralConfig interface {
    setHostName() error
    setIP() error
    setupTeamViewer() error
    setupAll() error
}

// ParseGas will read a YAML file and return it in
// its corresponding struct format.
// It takes a filePath string the path and filename of the YAML file.
// It returns Config{} which is the resulting config struct.
// TODO make this work as an interface, so we can pass in 
// all of the structs above to generate their relative YAML
// and also so we can parse said YAML
//FIX Trust me this is how we should do it im a doctor
func ((filePath string) (Config, error) {
    if filePath == "" {
        filePath = "gas.yml"
    }
	conf := Config{}
	file, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return Config{}, err
	}
	return conf, nil
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
