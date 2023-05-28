// Package core Module file
//
// This file contains everything directly related to handling a Module model.
package core

import (
	"errors"
	"gopkg.in/yaml.v3"
	"log"
)

const ModuleNotFoundError string = "no module found"

// Module
// Any struct that implements these methods can be considered a module.
type Module interface {
	SetConfig(*ModuleConfig)
	Config() ModuleConfig
	Setup() error
	TearDown() error
	Update() error
}

type ModuleConfig struct {
	Priority int
	Enabled  bool
	Settings interface{}
}

// Registry
// This registry contains Modules.
// The Modules take a string and map it to a Module.
type Registry struct {
	Modules map[string]Module
}

// ModuleRegistry
// This is the ModuleRegistry for the running application.
var ModuleRegistry = &Registry{
	make(map[string]Module),
}

func NewModuleRegistry() *Registry {
	return &Registry{
		make(map[string]Module),
	}
}

// Get
// This gets a module from an existing registry.
func (mr *Registry) Get(mod string) (Module, error) {
	found := mr.Modules[mod]
	if found != nil {
		return found, nil
	}
	return nil, errors.New(ModuleNotFoundError)
}

// Register
// This adds a new module to an existing registry.
func (mr *Registry) Register(name string, mod Module) {
	mr.Modules[name] = mod
}

// RunSetup
// Runs the Setup method on each Registry.Modules
func (mr *Registry) RunSetup() (err error) {
	err = ReadConfigFromFile("")
	// TODO handle priority of ModuleRegistry
	for _, module := range mr.Modules {
		log.Println(module.Config())
		if err != nil {
			return
		}
	}
	return
}

// RunUpdate
// Runs the Setup method on each Registry.Modules
func (mr *Registry) RunUpdate() (err error) {
	// TODO handle priority of ModuleRegistry
	for _, module := range mr.Modules {
		err = module.Update()
		if err != nil {
			return
		}
	}
	return
}

// setCurrent
// Sets the config for each module in the repository from the settingsYAML.
func (mr *Registry) setCurrent(settingsYAML []byte) error {
	// Unmarshal the YAML data into the ModuleSettings map.
	err := yaml.Unmarshal(settingsYAML, &moduleSettings)
	if err != nil {
		return err
	}

	// For each module in the registry, retrieve its settings from
	// the ModuleSettings map and set them.
	for moduleName, module := range mr.Modules {
		if settings, ok := moduleSettings[moduleName]; ok {
			info := module.Config()
			info.Settings = settings
			module.SetConfig(&info)
		}
	}

	return nil
}
