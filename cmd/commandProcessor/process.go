// commandprocessor is the first logical gate our tool encounters
package commandProcessor

import (
	"os"
	"sync"

	"github.com/Linkinlog/gasible/cmd/installer"
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

// Start the machine, handle which services to set up.
func initProcess() error {
	// Create a config struct and fill it from the config file.
	conf := models.Config{}.FillFromFile("")
	// Create a waitgroup so we can run all services at once.
	var wg sync.WaitGroup

	if conf.ServicesConfig.Installer {
		wg.Add(1)
		go installer.Installer(&conf.PackageInstallerConfig, &wg)
	}
	if conf.ServicesConfig.Ssh {
		// TODO
	}
	if conf.ServicesConfig.Teamviewer {
		// TODO
	}
	wg.Wait()
	return nil
}
