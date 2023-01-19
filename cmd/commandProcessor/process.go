// commandprocessor is the first logical gate our tool encounters
package commandProcessor

import (
	"sync"

	"github.com/Linkinlog/gasible/cmd/installer"
	"github.com/Linkinlog/gasible/internal/models"
)

// Start the machine, handle which services to set up.
func InitProcess(conf *models.Config) error {
	// Create a waitgroup so we can run all services at once.
	var wg sync.WaitGroup
	var errChan chan error

	if conf.ServicesConfig.Installer {
		wg.Add(1)
		go func() {
			err := installer.Run(&conf.PackageInstallerConfig, &conf.GlobalOpts, &wg)
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
	wg.Wait()
	if err := <-errChan; err != nil {
		return err
	}
	return nil
}
