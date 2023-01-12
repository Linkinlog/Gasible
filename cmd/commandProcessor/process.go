// commandprocessor is the first logical gate our tool encounters
package commandProcessor

import (
	"os"
	"runtime"

	"github.com/Linkinlog/gasible/cmd/osHandler"
	"github.com/Linkinlog/gasible/cmd/yamlParser"
	"github.com/Linkinlog/gasible/internal/models"
)

// ProcessCommand will start everything,
// handling the flags/args and processing.
func ProcessCommand() error {
	if len(os.Args) > 1 && os.Args[1] == "generate" {
		return yamlParser.CreateDefaults()
	}
	return initProcess()
}

func initProcess() error {
    // Grab system running this software
	system, err := osHandler.StringToSystem(runtime.GOOS)
	if err != nil {
		return err
	}
    // Create a config struct for us to fill
    conf := yamlParser.Config{}
    conf, err = conf.Parse()
	if err != nil { // Handle errors filling (?) maybe
		return err
	}
    // 
	if conf.ServicesConfig.Installer {
		// TODO Make go routine maybe
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
