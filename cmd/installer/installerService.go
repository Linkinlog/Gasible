package installer

import (
	"github.com/Linkinlog/gasible/internal/models"
)

// The relative options that we
// will need to pass into Run.
type InstallerOpts struct {
	NoOp bool
	Os   *models.System
}

// Install the packages listed in the
// packages section of the YAML file.
func (opts *InstallerOpts) Run(c *models.PackageInstallerConfig) ([]byte, error) {
	// Format our command
	command := c.GetCmd()

	return opts.Os.Exec(command)
}
