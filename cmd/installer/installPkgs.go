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
func Installer(c *models.PackageInstallerConfig, w *sync.WaitGroup) error {
	defer w.Done()
	// Verify the OS is supported.
	system := osHandler.GetCurrentSystem()
	// Validate the package manager and get its root path.
	pm := c.CheckPMAndReturnPath()
	// Turn our slice of packages into a single string.
	packages := strings.Join(c.Packages, " ")
	// Format all of the above into a string.
	command := strings.Join([]string{"echo", pm, c.Args, packages}, " ")

	// Mac / Linux have the same path for env
	if system.Name == "linux" || system.Name == "mac" {
		fmt.Println("Opening installLog.txt for logging..")
		// Open our log file for writing
		f, err := os.OpenFile("installLog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			m := fmt.Sprint("Error opening file:", err)
			panic(m)
		}
		defer f.Close()
		// Execute the command and prepare the output
		fmt.Println("Installing packages...")
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
