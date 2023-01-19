package installer

import (
	"errors"
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
func Installer(c *models.PackageInstallerConfig, g *models.GlobalOpts, w *sync.WaitGroup) error {
	defer w.Done()
	// Verify the OS is supported.
	system := osHandler.GetCurrentSystem()
	// Format our command
	command := setCmd(c)
	// If we are running in NoOp mode, only echo the command
	if g.NoOp {
		command = "echo Would have ran: " + command
	}
	// Mac / Linux have the same path for env
	if system.Name == "linux" || system.Name == "mac" {
		fmt.Println("Installing packages...")
		// Execute the command and prepare the output
		out, err := exec.Command("/usr/bin/env", "sh", "-c", command).Output()
		if err != nil {
			// Write the output of the error
			if err := log([]byte(err.Error())); err != nil {
				return err
			}
		}
		if err := log(out); err != nil {
			return err
		}

	} else if system.Name == "windows" {
		// TODO use the struct and make an interface to call this
		fmt.Println("TODO")
	}
	fmt.Println("Finished installing! Check the log for details.")
	return nil
}

func setCmd(c *models.PackageInstallerConfig) string {
	// Validate the package manager and get its root path.
	pm := c.CheckPMAndReturnPath()
	// Turn our slice of packages into a single string.
	packages := strings.Join(c.Packages, " ")
	// Format all of the above into a string.
	command := strings.Join([]string{pm, c.Args, packages}, " ")
	return command
}

func log(data []byte) error {
	// TODO Make this optional
	fmt.Println("Opening installLog.txt for logging..")
	// Open our log file for writing
	f, err := os.OpenFile("installLog.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		m := fmt.Sprint("Error opening file:", err)
		return errors.New(m)
	}
	defer f.Close()
	// Write the output
	_, err = f.Write(data)
	if err != nil {
		m := fmt.Sprint("Error writing to file:", err)
		return errors.New(m)
	}
	return nil
}
