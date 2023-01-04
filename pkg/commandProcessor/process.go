// commandprocessor is the first logical gate our tool encounters
package commandprocessor

import (
	"os"
	"runtime"

	"github.com/Linkinlog/gasible/pkg/osHandler"
	"github.com/Linkinlog/gasible/pkg/installer"
	yamlparser "github.com/Linkinlog/gasible/pkg/yamlParser"
)

// ProcessCommand will start everything,
// handling the flags/args and processing.
func ProcessCommand() error {
	if len(os.Args) > 1 && os.Args[1] == "generate" {
		return yamlparser.CreateDefaults()
	}
	return initProcess()
}

func initProcess() error {
	system, err := osHandler.StringToSystem(runtime.GOOS)
	if err != nil {
		return err
	}
	conf, err := yamlparser.ParseGas()
	if err != nil {
		return err
	}
    if conf.ServicesConfig.Installer {
        // TODO Make go routine
        installer.Installer(system, conf.PackageInstallerConfig)
    }
    if conf.ServicesConfig.Ssh {
        // TODO
    }
    if conf.ServicesConfig.Teamviewer {
        // TODO
    }
	return err
}
