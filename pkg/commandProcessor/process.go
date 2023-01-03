// commandprocessor is the first logical gate our tool encounters
package commandprocessor

import (
	"os"
	"runtime"

	"github.com/Linkinlog/gasible/pkg/installer"
	"github.com/Linkinlog/gasible/pkg/osHandler"
	yamlparser "github.com/Linkinlog/gasible/pkg/yamlParser"
)

// ProcessCommand will start everything,
// handling the flags/args and processing.
func ProcessCommand() error {
	system, err := osHandler.StringToSystem(runtime.GOOS)
	if err != nil {
		return err
	}
	if len(os.Args) > 1 && os.Args[1] == "generate" {
		return yamlparser.CreateDefaults()
	} else {
		return installer.Installer(system)
	}
}
