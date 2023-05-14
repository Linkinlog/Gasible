package core_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/Linkinlog/gasible/internal/core"
)

type InstallerTestCase struct {
	system *core.System
	Name   string
}

type MockCommandRunner struct{}

func (m MockCommandRunner) Command(name string, arg ...string) *exec.Cmd {
	cmdString := []string{"-test.run=TestSystemMock", "--", name}
	cmdString = append(cmdString, arg...)
	cmd := exec.Command(os.Args[0], cmdString...)
	cmd.Env = []string{"RUN_SYSTEM_MOCK=1"}
	return cmd
}

func TestSystemMock(t *testing.T) {
	if os.Getenv("RUN_SYSTEM_MOCK") != "1" {
		return
	}
	_, err := fmt.Fprint(os.Stdout, "mocking passed")
	if err != nil {
		t.Fatal("Failed mocking")
	}
	os.Exit(0)
}

func TestSystemExec(t *testing.T) {
	mockTestRunner := MockCommandRunner{}
	system := core.System{
		Name:   "TestOs",
		Runner: mockTestRunner,
	}
	testCase := InstallerTestCase{
		&system,
		"TestInstallerWithDefaults",
	}
	t.Run(testCase.Name, func(t *testing.T) {
		out, err := testCase.system.Exec("")
		expectedRes := "mocking passed"
		if err != nil {
			t.Fatalf("Failed installer test, err: %s", err.Error())
		}
		if string(out) != expectedRes {
			t.Fatalf("Failed installer mocking, wanted %s, got %s", expectedRes, string(out))
		}
	})
}
