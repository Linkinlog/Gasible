// Gasible
// Written by Dahlton

// Gasible aims to assist in automating setting up
// a local development environment by making or using
// a bare Git repository and Git submodules to manage
// dotfiles in an organized and effective manner

// Check out https://github.com/Linkinlog/gasible for more
package main

import (
	commandprocessor "github.com/Linkinlog/gasible/pkg/commandProcessor"
)

// main starts everything off, starting with our command processor.
func main() {
	commandprocessor.ProcessCommand()
}
