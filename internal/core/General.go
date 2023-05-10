package core

// The credentials to the teamviewer account.
// TODO Find a better way to handle this
type TVCreds struct {
	User string `yaml:"user,omitempty"`
	Pass string `yaml:"pass,omitempty"`
}

// GeneralConfig holds general configurations related to system setup.
type GeneralConfig struct {
	Hostname        string  `yaml:"hostname,omitempty"`
	StaticIP        string  `yaml:"staticIP,omitempty"`
	SubnetMask      string  `yaml:"subnetMask,omitempty"`
	TeamViewerCreds TVCreds `yaml:"TeamViewerCreds,omitempty"`
}

// Default creates a GeneralConfig object with default values.
func (GeneralConfig) Default() *GeneralConfig {
	return &GeneralConfig{
		Hostname:   "development-station",
		StaticIP:   "192.168.4.20",
		SubnetMask: "255.255.255.0",
		TeamViewerCreds: TVCreds{
			User: "username",
			Pass: "password",
		},
	}
}
