package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/Linkinlog/gasible/cmd/osHandler"
	"github.com/Linkinlog/gasible/internal/models"
)

// Install the packages listed in the
// packages section of the YAML file.
func Installer(c *models.Config, w *sync.WaitGroup) error {
	defer w.Done()
	// Verify the OS is supported.
	system := osHandler.GetCurrentSystem()
	// Validate the package manager and get its root path.
	pm := c.PackageInstallerConfig.CheckPMAndReturnPath()
	// Grab our install config from the config
	config := c.PackageInstallerConfig
	// Turn our slice of packages into a single string.
	packages := strings.Join(config.Packages, " ")
	// Format all of the above into a string.
	command := strings.Join([]string{pm, config.Args, packages}, " ")

	// Mac / Linux have the same path for env
	if system.Name == "linux" || system.Name == "mac" {
		fmt.Println("Installing packages...")
		// Open our log file for writing
		f, err := os.OpenFile("installLog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			m := fmt.Sprint("Error opening file:", err)
			panic(m)
		}
		defer f.Close()
		// Execute the command and prepare the output
		out, err := exec.Command("/usr/bin/env", "sh", "-c", command).Output()
		if err != nil {
			// Write the output
			_, err := f.Write([]byte(err.Error()))
			if err != nil {
				panic(err)
			}
		}

		// Write the output
		_, err = f.Write(out)
		if err != nil {
			panic(err)
		}
	} else if system.Name == "windows" {
		// TODO
	}
	fmt.Println("Finished installing packages:", packages)
	return nil
}
