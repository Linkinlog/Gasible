package modules

import (
	"github.com/Linkinlog/gasible/internal"
	"github.com/Linkinlog/gasible/internal/app"
	"io"
	"log"
	"os/exec"
	"runtime"
)

// init
// This should really just handle registering the module in the registry.
func init() {
	ToBeRegistered = append(ToBeRegistered, &SysCall{
		name:    "SysCall",
		Enabled: true,
		Settings: sysCallSettings{
			CurrentOS: runtime.GOOS,
		},
		sysCommand: cmdRunner{},
	})
}

// SysCall implements the module interface, so we can execute system commands.
type SysCall struct {
	name        string
	Enabled     bool
	Settings    sysCallSettings
	application *app.App
	sysCommand
}

// sysCommand is what executes commands on the system.
type sysCommand interface {
	Exec(command string, args []string, sudo bool) ([]byte, error)
	ExecWithInput(command string, args []string, stdinInput string, sudo bool) ([]byte, error)
}

// sysCallSettings allows us to keep track of the running OS.
type sysCallSettings struct {
	CurrentOS string `yaml:"-"`
}

// ParseConfig does nothing as there is no config for this module.
func (s *SysCall) ParseConfig(_ map[string]interface{}) error {
	return nil // No config for this module
}

// Config returns the shallow-copied module config from our module's config.
func (s *SysCall) Config() app.ModuleConfig {
	return app.ModuleConfig{
		Enabled:  s.Enabled,
		Settings: s.Settings,
	}
}

// GetName returns the name field of the SysCall struct.
func (s *SysCall) GetName() string {
	return s.name
}

// Setup runs nothing as we have nothing to do for this module.
func (s *SysCall) Setup() error {
	return nil // no setup required
}

// TearDown runs nothing as we have nothing to do for this module.
func (s *SysCall) TearDown() error {
	return nil // no teardown required
}

// Update runs nothing as we have nothing to do for this module.
func (s *SysCall) Update() error {
	return nil
}

// SetApp sets the application field as the app that is passed in.
func (s *SysCall) SetApp(app *app.App) {
	s.application = app
}

// cmdRunner implements SysCommand as a base command runner.
type cmdRunner struct{}

// Exec for when we need to execute a command on the host system.
func (r cmdRunner) Exec(command string, args []string, sudo bool) ([]byte, error) {
	if sudo {
		args = append([]string{command}, args...)
		command = "sudo"
	}
	execCmd := exec.Command(command, args...)
	log.Println("Executing: " + execCmd.String())
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return nil, internal.ErrorAs("cmdRunner.Exec", err)
	}
	return output, nil
}

// ExecWithInput for instances where we need to simulate piping something into a command.
func (r cmdRunner) ExecWithInput(command string, args []string, stdinInput string, sudo bool) ([]byte, error) {
	if sudo {
		args = append([]string{command}, args...)
		command = "sudo"
	}
	execCmd := exec.Command(command, args...)

	stdin, err := execCmd.StdinPipe()
	if err != nil {
		return []byte{}, internal.ErrorAs("ExecWithInput", err)
	}

	_, writeErr := io.WriteString(stdin, stdinInput)
	if writeErr != nil {
		return []byte{}, internal.ErrorAs("ExecWithInput", writeErr)
	}

	if closeErr := stdin.Close(); closeErr != nil {
		return []byte{}, internal.ErrorAs("ExecWithInput", closeErr)
	}

	log.Println("Executing: " + execCmd.String())
	return execCmd.CombinedOutput()
}
