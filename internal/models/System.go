// Systems file
//
// This file contains structs relative to
// the operating system that we are setting up.
// As well as functions to validate and set the target system.
package models

import (
	"errors"
	"log"
	"os/exec"
	"runtime"
)

// System contains the name of the operating system and
// the configs we want to configure for said OS.
type System struct {
	Name    string
	Cmd     Cmd
	Configs []config
}

// config contains the name of the configuration and
// the path we should put the config files / clone the repo to.
type config struct {
	Name string
	Path string
}

// Fields relative to the command we will execute.
type Cmd struct {
	Exec string
	Args []string
	Env  []string
}

// Executes the command string, or echo's out the command it would have ran.
func (os System) Exec(noop bool, command string) error {
	// Set up the command and handle noop
	var execCmd *exec.Cmd
	if noop {
		args := append([]string{"Would have ran: ", os.Cmd.Exec}, os.Cmd.Args...)
		execCmd = exec.Command("echo", args...)
	} else {
		execCmd = exec.Command(os.Cmd.Exec, os.Cmd.Args...)
	}
	execCmd.Args = append(execCmd.Args, command)
	out, err := execCmd.Output()
	if err != nil {
		// Write the output of the error to the program log
        return err
	} else {
		log.Print(string(out))
	}
	return nil
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
	//"windows": windowsSystem,
	"linux": linuxSystem,
	//"darwin":  macSystem,
}

// Default values for a linux OS.
var linuxSystem = System{
	"linux",
	Cmd{
		"/usr/bin/env", []string{"sh -c"}, []string{}},
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

// TODO: create Win/Mac
