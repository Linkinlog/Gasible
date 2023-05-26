package core

// PackageManagerConfig holds all fields
// relative to the package installer service.
type PackageManagerConfig struct {
	Manager  string   `yaml:"pkg-manager-command,omitempty"`
	Packages []string `yaml:"packages"`
}

// Default Populate the struct with the default config for the package installer.
func (*PackageManagerConfig) Default() *PackageManagerConfig {
	return &PackageManagerConfig{
		Manager: "apt",
		Packages: []string{
			"neovim",
		},
	}
}
