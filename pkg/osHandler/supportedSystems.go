// SupportSystems file
//
// This file contains structs relative to
// the operating system that we are setting up.
// As well as functions to validate and set the target system.
package osHandler

import (
	"errors"
)

// config contains the name of the configuration and
// the path we should put the config files / clone the repo to.
type config struct {
	name string
	path string
}

// os contains the name of the operating system and
// the configs we want to configure for said OS.
type System struct {
	name       string
	pkgManager string
	configs    []config
}

// stringToSystem will take a system as a string value
// and return it in a struct format as long as it is supported.
// It will return an error if it is not supported.
func StringToSystem(s string) (System, error) {
	switch s {
	case "win", "linux", "mac":
		return supportedSystemsMap[s], nil
	default:
		return System{}, errors.New("Unsupported system")
	}
}

// supportedSystemsMap is a map so we can access
// a System struct by its string value.
var supportedSystemsMap = map[string]System{
	"win":   windowsSystem,
	"linux": linuxSystem,
	"mac":   macSystem,
}

// Default values for a windows OS.
var windowsSystem = System{
	"win",
	"winget",
	[]config{
		{
			"neovim",
			"TBD",
		},
		{
			"tmux",
			"TBD",
		},
		{
			"zsh",
			"TBD",
		},
	},
}

// Default values for a linux OS.
var linuxSystem = System{
	"linux",
	"dnf",
	[]config{
		{
			"neovim",
			".config/nvim/",
		},
		{
			"tmux",
			".config/tmux/",
		},
		{
			"zsh",
			".config/zsh/",
		},
	},
}

// Default values for a mac OS.
var macSystem = System{
	"mac",
	"brew",
	[]config{
		{
			"neovim",
			".config/nvim/",
		},
		{
			"tmux",
			".config/tmux/",
		},
		{
			"zsh",
			".config/zsh/",
		},
	},
}
