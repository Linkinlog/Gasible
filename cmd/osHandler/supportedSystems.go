// SupportSystems file
//
// This file contains structs relative to
// the operating system that we are setting up.
// As well as functions to validate and set the target system.
package osHandler

import (
	"errors"
	"runtime"
)

// config contains the name of the configuration and
// the path we should put the config files / clone the repo to.
type config struct {
	Name string
	Path string
}

// os contains the name of the operating system and
// the configs we want to configure for said OS.
type System struct {
	Name    string
	Configs []config
}

// GetCurrentSystem returns the system struct
// relative to the current system at runtime.
func GetCurrentSystem() *System {
	// Grab system running this software
	system, err := StringToSystem(runtime.GOOS)
	if err != nil {
		panic(err)
	}
	return &system
}

// stringToSystem will take a system as a string value
// and return it in a struct format as long as it is supported.
// It will return an error if it is not supported.
func StringToSystem(s string) (System, error) {
	if system, exists := supportedSystemsMap[s]; exists {
		return system, nil
	} else {
		return System{}, errors.New("Unsupported system")
	}
}

// supportedSystemsMap is a map so we can access
// a System struct by its string value.
var supportedSystemsMap = map[string]System{
	"windows": windowsSystem,
	"linux":   linuxSystem,
	"darwin":  macSystem,
}

// Default values for a windows OS.
var windowsSystem = System{
	"windows",
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
