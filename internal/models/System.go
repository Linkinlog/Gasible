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
func (os System) Exec(noop bool, command string) ([]byte, error) {
	// Set up the command and handle noop
	var execCmd *exec.Cmd
	if noop {
		args := append([]string{"Would have ran: ", os.Cmd.Exec}, os.Cmd.Args...)
		execCmd = exec.Command("echo", args...)
	} else {
		execCmd = exec.Command(os.Cmd.Exec, os.Cmd.Args...)
	}
	execCmd.Args = append(execCmd.Args, command)
	log.Print("Beginning package installation! Please wait...")
	out, err := execCmd.CombinedOutput()
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}
