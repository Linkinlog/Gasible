package installer

import (
	"github.com/Linkinlog/gasible/internal/core"
)

// Opts InstallerOpts is the relative options that we
// will need to pass into Run.
type Opts struct {
	NoOp bool
	Os   *core.System
}

// Run installs the packages listed in the
// packages section of the YAML file.
func (opts *Opts) Run(c *core.PackageManagerConfig) ([]byte, error) {
	// Format our command
	command, err := c.GetInstallCommand()
	if err != nil {
		return []byte{}, err
	}

	return opts.Os.Exec(command)
}
