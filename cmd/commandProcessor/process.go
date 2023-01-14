// commandprocessor is the first logical gate our tool encounters
package commandProcessor

import (
	"os"

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

func initProcess() error {
	// Create a config struct and fill it from the config file
	conf := models.Config{}.FillFromFile("")

	if conf.ServicesConfig.Installer {
		// TODO Make go routine maybe
		installer.Installer(conf)
	}
	if conf.ServicesConfig.Ssh {
		// TODO
	}
	if conf.ServicesConfig.Teamviewer {
		// TODO
	}
	return nil
}
