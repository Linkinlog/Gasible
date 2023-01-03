// commandprocessor is the first logical gate our tool encounters
package commandprocessor

import (
	"os"
	"runtime"

	"github.com/Linkinlog/gasible/pkg/installer"
	"github.com/Linkinlog/gasible/pkg/osHandler"
)

// ProcessCommand will start everything,
// handling the flags/args and processing.
func ProcessCommand() error {
	system, err := osHandler.StringToSystem(runtime.GOOS)
	if err != nil {
		return err
	}
		return installer.Installer(system)
}
