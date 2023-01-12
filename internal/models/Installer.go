package models

type PackageInstallerConfig struct {
	Manager  string   `yaml:"pkg-manager,omitempty"`
	Packages []string `yaml:"packages,omitempty"`
}

func (pkgInstallConf PackageInstallerConfig) Default() *PackageInstallerConfig{
	return &PackageInstallerConfig{
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
}

