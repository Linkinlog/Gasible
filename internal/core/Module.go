// Package core Module file
//
// This file contains everything directly related to handling a Module model.
package core

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
)

var ModuleNotFoundError = errors.New("no module found")
var dependencyGraphCycleError = errors.New("there is a cycle in the module dependencies")

type moduleAction func(Module) error

// Module
// Any struct that implements these methods can be considered a module.
type Module interface {
	ParseConfig(map[string]interface{}) error
	Config() ModuleConfig
	Setup() error
	TearDown() error
	Update() error
	GetDeps() []string
}

// ModuleConfig
// General items we may need to track for each module.
type ModuleConfig struct {
	Enabled  bool        `yaml:"enabled"`
	Settings interface{} `yaml:"settings"`
}

// Registry
// This registry contains Modules.
// The Modules take a string and map it to a Module.
type Registry struct {
	Modules      map[string]Module
	Dependencies map[string][]string
}

func NewModuleRegistry() *Registry {
	return &Registry{
		Modules:      make(map[string]Module),
		Dependencies: make(map[string][]string),
	}
}

// ModuleRegistry is the ModuleRegistry for the running application.
var ModuleRegistry = NewModuleRegistry()

// Get a module from an existing registry.
func (mr *Registry) Get(mod string) (Module, error) {
	found := mr.Modules[mod]
	if found != nil {
		return found, nil
	}
	return nil, ModuleNotFoundError
}

// Register adds a new module to an existing registry.
func (mr *Registry) Register(name string, mod Module) {
	mr.Modules[name] = mod
	mr.Dependencies[name] = mod.GetDeps()
}

// RunSetup Runs the Setup method on each Registry.Modules
func (mr *Registry) RunSetup() (err error) {
	// TODO handle priority of ModuleRegistry
	for _, module := range mr.Modules {
		if !module.Config().Enabled {
			return nil
		}
		err = module.Setup()
		if err != nil {
			return
		}
	}
	return
}

// RunUpdate runs the Setup method on each Registry.Modules
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

// RunTeardown runs the Setup method on each Registry.Modules
func (mr *Registry) RunTeardown() (err error) {
	// TODO handle priority of ModuleRegistry
	for _, module := range mr.Modules {
		err = module.TearDown()
		if err != nil {
			return
		}
	}
	return
}

// setCurrent sets the Config for each module in the repository from the settingsYAML.
func (mr *Registry) setCurrent(settingsYAML []byte) error {
	// Unmarshal the YAML data into the ModuleSettings map.
	err := yaml.Unmarshal(settingsYAML, &moduleSettingsMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshall yaml: %w", err)
	}

	// For each module in the registry, retrieve its settings from
	// the ModuleSettings map and set them.
	for moduleName, module := range mr.Modules {
		if rawSettings, ok := moduleSettingsMap[moduleName]; ok {
			err = module.ParseConfig(rawSettings.(map[string]interface{}))
			if err != nil {
				return fmt.Errorf("parseConfig failed parsing: %w", err)
			}
		}
	}

	return nil
}
