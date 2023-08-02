// Gasible
// Written by Dahlton

// Gasible aims to assist in automating setting up
// a local development environment by making or using
// a bare Git repository and Git submodules to manage
// dotfiles in an organized and effective manner

// Check out https://github.com/Linkinlog/gasible for more
package main

import (
	"github.com/Linkinlog/gasible/cmd"
	"github.com/Linkinlog/gasible/internal/app"
	"github.com/Linkinlog/gasible/modules"
	_ "github.com/Linkinlog/gasible/modules"
	"log"
)

// main starts everything off, now handled by Cobra.
func main() {
	newApp := app.New()
	for _, module := range modules.ToBeRegistered {
		module.SetApp(newApp)
		newApp.ModuleRegistry.Register(module)
	}
	err := cmd.ExecuteApplication(newApp)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
