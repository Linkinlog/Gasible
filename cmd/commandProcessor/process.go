// commandprocessor is the first logical gate our tool encounters
package commandProcessor

import (
	"log"
	"os"
	"sync"

	"github.com/Linkinlog/gasible/cmd/gitService"
	"github.com/Linkinlog/gasible/cmd/installer"
	"github.com/Linkinlog/gasible/internal/models"
)

// Start the machine, handle which services to set up.
func InitProcess(conf *models.Config) error {
	// Create a waitgroup so we can run all services at once.
	var wg sync.WaitGroup
	// Open the log file for writing
	logFile := "app.log"
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.LstdFlags)

	if conf.ServicesConfig.Installer {
		wg.Add(1)
		go func() {
			defer wg.Done()
			opts := installer.InstallerOpts{
				NoOp: conf.GlobalOpts.NoOp,
				Os:   models.GetCurrentSystem(),
			}
			out, err := opts.Run(&conf.PackageInstallerConfig)
			if err != nil {
				log.Fatal(err)
				logger.Fatal(err)
			} else if out != nil {
				logger.Println(string(out))
			}
		}()
	}
	if conf.GitServiceConfig.Enabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			opts := gitService.Opts{
				NoOp: conf.GlobalOpts.NoOp,
			}
			out, err := opts.Run(&conf.GitServiceConfig)
			if err != nil {
				logger.Println("Error: " + err.Error())
			}
			for _, e := range out {
				if len(e) > 0 {
					logger.Print(string(e))
				}
			}
		}()
	}

	// if conf.ServicesConfig.Ssh {
	// TODO
	// }
	// if conf.ServicesConfig.Teamviewer {
	// TODO
	// }
	// wait for all the goroutines to complete
	wg.Wait()

	return nil
}
