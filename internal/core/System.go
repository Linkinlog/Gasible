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

type CmdRunner interface {
	Command(name string, arg ...string) *exec.Cmd
}

type SudoRunner struct{}

func (r SudoRunner) Command(name string, arg ...string) *exec.Cmd {
	_, err := exec.LookPath(name)
	if err != nil {
		return nil
	}
	arg = append([]string{"sudo", name}, arg...)
	stringArg := strings.Join(arg, " ")
	return exec.Command("/bin/sh", "-c", stringArg)
}

// System contains the name of the operating system and
// the configs we want to configure for said OS.
type System struct {
	Name   string
	Runner CmdRunner
}

// Exec executes the command string.
func (os System) Exec(command string, args []string) ([]byte, error) {
	execCmd := os.Runner.Command(command, args...)
	out, err := execCmd.CombinedOutput()
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}
