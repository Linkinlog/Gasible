package core

// Information relatie to the services we will configure.
type ServicesConfig struct {
	Installer  bool `yaml:"installer,omitempty"`
	Teamviewer bool `yaml:"teamviewer,omitempty"`
	Ssh        bool `yaml:"ssh,omitempty"`
	Git        bool `yaml:"git,omitempty"`
}

// Sets the defaults.
func (ServicesConfig) Default() *ServicesConfig {
	return &ServicesConfig{
		Installer:  true,
		Teamviewer: true,
		Ssh:        true,
		Git:        true,
	}
}
