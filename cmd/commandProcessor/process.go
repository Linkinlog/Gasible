// commandprocessor is the first logical gate our tool encounters
package commandProcessor

import (
	"log"
	"sync"

	"github.com/Linkinlog/gasible/cmd/installer"
	"github.com/Linkinlog/gasible/internal/models"
)

// Start the machine, handle which services to set up.
func InitProcess(conf *models.Config) error {
	// Create a waitgroup so we can run all services at once.
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	if conf.ServicesConfig.Installer {
		wg.Add(1)
		go func() {
			defer wg.Done()
			opts := installer.InstallerOpts{
				NoOp: conf.GlobalOpts.NoOp,
				Os:   models.GetCurrentSystem(),
			}
			err := opts.Run(&conf.PackageInstallerConfig)
			if err != nil {
				errChan <- err
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

	// check if there were any errors
	select {
	case err := <-errChan:
		log.Fatalf("errorChan: %v", err)
		close(errChan)
	default:
		log.Println("all goroutines completed successfully")
	}
	return nil
}
