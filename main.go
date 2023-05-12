// Gasible
// Written by Dahlton

// Gasible aims to assist in automating setting up
// a local development environment by making or using
// a bare Git repository and Git submodules to manage
// dotfiles in an organized and effective manner

// Check out https://github.com/Linkinlog/gasible for more
package main

import (
	_ "embed"
	"github.com/Linkinlog/gasible/cmd"
	"github.com/Linkinlog/gasible/internal/core"
	_ "github.com/Linkinlog/gasible/modules"
)

// Embedding the config.yml so we can include it in the binary
// Allows us to have a portable executable
// TODO handle all yariants of yaml files
//
//go:embed config.yml
var configFile string

// This is the ConfigModel for the running application.
var Config = core.NewConfigFromFile(configFile)

// main starts everything off, now handled by Cobra.
func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
