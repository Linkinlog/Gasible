package installer_test

import (
	"testing"

	"github.com/Linkinlog/gasible/internal/models"
)

func TestInstallerWithDefaults(t *testing.T) {
	type InstallerTestCase struct {
		s    *models.System
		c    string
		Name string
	}
	testCase := InstallerTestCase{
		models.GetCurrentSystem(),
		models.PackageInstallerConfig{}.Default().GetCmd(),
		"TestInstallerWithDefaults",
	}
	t.Run(testCase.Name, func(t *testing.T) {
		if err := testCase.s.Exec(true, testCase.c); err != nil {
			t.Fatalf("Failed installer test, err: %s", err.Error())
		}
	})
}
