// Package app Systems file
//
// This file contains structs relative to
// the operating system that we are setting up.
// As well as functions to validate and set the target system.
package app

import (
	"io"
	"os/exec"
	"runtime"
)

const sudo = "sudo"

// The CurrentSystem contains the name of the operating system and
// the configs we want to configure, for said OS.
type CurrentSystem struct {
	Name string
}

func NewSystem() *CurrentSystem {
	return &CurrentSystem{
		Name: runtime.GOOS,
	}
}

// cmdExecutor
// Allows us to mock running commands for testing.
type cmdExecutor interface {
	Command(name string, arg ...string) (*exec.Cmd, error)
}

func (sys *CurrentSystem) ExecCombinedOutput(runner cmdExecutor, command string, args []string) ([]byte, error) {
	execCmd, err := runner.Command(command, args...)
	if err != nil {
		return []byte{}, ErrorAs("ExecCombinedOutput", err)
	}
	return execCmd.CombinedOutput()
}

func (sys *CurrentSystem) ExecRun(runner cmdExecutor, command string, args []string) error {
	execCmd, err := runner.Command(command, args...)
	if err != nil {
		return ErrorAs("ExecRun", err)
	}
	startErr := execCmd.Run()
	if startErr != nil {
		return ErrorAs("ExecRun", startErr)
	}
	return nil
}

func (sys *CurrentSystem) ExecCombinedWithInput(runner cmdExecutor, command string, args []string, stdinInput string) ([]byte, error) {
	execCmd, err := runner.Command(command, args...)
	if err != nil {
		return []byte{}, ErrorAs("ExecCombinedWithInput", err)
	}
	stdin, err := execCmd.StdinPipe()
	if err != nil {
		return []byte{}, ErrorAs("ExecCombinedWithInput", err)
	}

	done := make(chan error, 1)
	go func() {
		defer close(done)
		_, writeErr := io.WriteString(stdin, stdinInput)
		if writeErr != nil {
			done <- ErrorAs("ExecCombinedWithInput", writeErr)
			return
		}
		if closeErr := stdin.Close(); closeErr != nil {
			done <- ErrorAs("ExecCombinedWithInput", closeErr)
			return
		}
	}()
	if chanErrs := <-done; chanErrs != nil {
		return []byte{}, chanErrs
	}

	return execCmd.CombinedOutput()
}

// SudoRunner
// cmdExecutor implementation that uses sudo.
type SudoRunner struct{}

func (r SudoRunner) Command(name string, arg ...string) (*exec.Cmd, error) {
	return createCmdWithPrefix([]string{sudo}, name, arg...)
}

// NormalRunner
// cmdExecutor implementation that doesn't use sudo.
type NormalRunner struct{}

func (r NormalRunner) Command(name string, arg ...string) (*exec.Cmd, error) {
	return createCmdWithPrefix([]string{}, name, arg...)
}

func createCmdWithPrefix(prefix []string, name string, arg ...string) (*exec.Cmd, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, ErrorAs("createCmdWithPrefix", err)
	}
	arg = append(prefix, append([]string{path}, arg...)...)
	return exec.Command(arg[0], arg[1:]...), nil
}
