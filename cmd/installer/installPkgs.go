package installer

import (
	"github.com/Linkinlog/gasible/cmd/osHandler"
	"github.com/Linkinlog/gasible/internal/models"
)

type InstallerOpts struct {
	NoOp bool
	Os   *osHandler.System
}

func Opts() *InstallerOpts {
	newInstallOpts := &InstallerOpts{}
	newInstallOpts.NoOp = false
	return newInstallOpts
}

// Install the packages listed in the
// packages section of the YAML file.
func (opts *InstallerOpts) Run(c *models.PackageInstallerConfig) error {
	// Format our command
	command := c.GetCmd()

	return opts.Os.Exec(opts.NoOp, command)
}
