// Package core Systems file
//
// This file contains structs relative to
// the operating system that we are setting up.
// As well as functions to validate and set the target system.
package core

import (
	"os/exec"
	"strings"
)

const shExec = "/bin/sh"
const shExecArg = "-c"
const sudo = "sudo"

// System contains the name of the operating system and
// the configs we want to configure for said OS.
type System struct {
	Name   string
	Runner cmdRunner
}

// cmdRunner
// Allows us to mock running commands for testing.
type cmdRunner interface {
	Command(name string, arg ...string) *exec.Cmd
}

// SudoRunner
// cmdRunner implementation that uses sudo.
type SudoRunner struct{}

// Exec executes the command string.
func (os System) Exec(command string, args []string) ([]byte, error) {
	execCmd := os.Runner.Command(command, args...)
	return execCmd.CombinedOutput()
}

// Command
// Uses sudo to run command name with args arg.
func (r SudoRunner) Command(name string, arg ...string) *exec.Cmd {
	_, err := exec.LookPath(name)
	if err != nil {
		return nil
	}
	arg = append([]string{sudo, name}, arg...)
	stringArg := strings.Join(arg, " ")
	return exec.Command(shExec, shExecArg, stringArg)
}
