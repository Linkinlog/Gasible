package modules

import (
	"github.com/Linkinlog/gasible/internal"
	"github.com/Linkinlog/gasible/internal/app"
	"io"
	"log"
	"os/exec"
	"runtime"
)

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

type SysCall struct {
	name        string
	Enabled     bool
	Settings    sysCallSettings
	application *app.App
	sysCommand
}

type sysCommand interface {
	Exec(command string, args []string, sudo bool) ([]byte, error)
	ExecWithInput(command string, args []string, stdinInput string, sudo bool) ([]byte, error)
}

type sysCallSettings struct {
	CurrentOS string `yaml:"-"`
}

func (s *SysCall) ParseConfig(_ map[string]interface{}) error {
	return nil // No config for this module
}

func (s *SysCall) Config() app.ModuleConfig {
	return app.ModuleConfig{
		Enabled:  s.Enabled,
		Settings: s.Settings,
	}
}

func (s *SysCall) GetName() string {
	return s.name
}

func (s *SysCall) Setup() error {
	return nil // no setup required
}

func (s *SysCall) TearDown() error {
	return nil // no teardown required
}

func (s *SysCall) Update() error {
	return nil // no Update required
}

func (s *SysCall) SetApp(app *app.App) {
	s.application = app
}

type cmdRunner struct{}

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
