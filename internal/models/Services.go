package models

type ServicesConfig struct {
	Installer  bool `yaml:"installer,omitempty"`
	Teamviewer bool `yaml:"teamviewer,omitempty"`
	Ssh        bool `yaml:"ssh,omitempty"`
	Git        bool `yaml:"git,omitempty"`
}

func (Services ServicesConfig) Default() *ServicesConfig {
	return &ServicesConfig{
		Installer:  true,
		Teamviewer: true,
		Ssh:        true,
		Git:        true,
	}
}
