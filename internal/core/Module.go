// Package core Module file
//
// This file contains everything directly related to handling a Module model.
package core

import (
	"errors"
)

const ModuleNotFoundError string = "no module found"

// Module Any struct that implements these methods can be considered a module.
type Module interface {
	Setup() error
	Update() error
}

// Registry This registry contains modules.
// The modules take a string and map it to a Module.
type Registry struct {
	modules map[string]Module
}

// ModuleRegistry This is the ModuleRegistry for the running application.
var ModuleRegistry = &Registry{
	make(map[string]Module),
}

func NewModuleRegistry() *Registry {
	return &Registry{
		make(map[string]Module),
	}
}

// Get This gets a module from an existing registry.
func (mr *Registry) Get(mod string) (Module, error) {
	found := mr.modules[mod]
	if found != nil {
		return found, nil
	}
	return nil, errors.New(ModuleNotFoundError)
}

// Register This adds a new module to an existing registry.
func (mr *Registry) Register(name string, mod Module) {
	mr.modules[name] = mod
}
