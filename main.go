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
	_ "github.com/Linkinlog/gasible/modules"
	"log"
)

// main starts everything off, now handled by Cobra.
func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
