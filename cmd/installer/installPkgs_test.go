package installer_test

import (
	"testing"

	"github.com/Linkinlog/gasible/cmd/installer"
	"github.com/Linkinlog/gasible/cmd/osHandler"
	"github.com/Linkinlog/gasible/internal/models"
)

func TestInstallerWithDefaults(t *testing.T) {
	type InstallerTestCase struct {
		s                *osHandler.System
		c                string
		noop             bool
		Name             string
	}
	testCase := InstallerTestCase{
		osHandler.GetCurrentSystem(),
		installer.GetCmd(models.PackageInstallerConfig{}.Default()),
		true,
		"TestInstallerWithDefaults",
	}
	t.Run(testCase.Name, func(t *testing.T) {
		if err := installer.Installer(testCase.s, testCase.c, testCase.noop); err != nil {
			t.Fatalf("Failed running installer.Installer, err: %s", err.Error())
		}
	})
}
