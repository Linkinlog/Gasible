// Systems file
//
// This file contains structs relative to
// the operating system that we are setting up.
// As well as functions to validate and set the target system.
package core

import (
	"log"
	"os/exec"
)

type CmdRunner interface {
	Command(name string, arg ...string) *exec.Cmd
}

type RealRunner struct{}

func (r RealRunner) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

// System contains the name of the operating system and
// the configs we want to configure for said OS.
type System struct {
	Name   string
	Runner CmdRunner
}

// Executes the command string, or echo's out the command it would have ran.
func (os System) Exec(command string, args ...string) ([]byte, error) {
	// Set up the command and handle noop
	execCmd := os.Runner.Command(command, args...)
	execCmd.Args = append(execCmd.Args, command)
	log.Print("Beginning package installation! Please wait...")
	out, err := execCmd.CombinedOutput()
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}
