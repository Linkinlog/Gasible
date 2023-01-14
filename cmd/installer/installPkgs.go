package installer

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/Linkinlog/gasible/cmd/osHandler"
	"github.com/Linkinlog/gasible/internal/models"
)

// Install the packages listed in the
// packages section of the YAML file.
func Installer(c *models.Config) error {
	// Verify the OS is supported
	system := osHandler.GetCurrentSystem()
    // Validate the package manager
    pm := c.PackageInstallerConfig.CheckPMAndReturnPath()
	config := c.PackageInstallerConfig
	packages := strings.Join(config.Packages, " ")
    command := strings.Join([]string{pm, config.Args, packages}, " ")

	if system.Name == "linux" || system.Name == "mac" {
        out, err := exec.Command("/usr/bin/env", "sh", "-c", command).Output()
        if err != nil {
            panic(err)
        }
        
	} else if system.Name == "windows" {
		exec.Command("/usr/bin/env sh", "-c", pm, config.Args, packages)
	}
	return errors.New("WIP")
}
