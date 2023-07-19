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

type moduleAction func(Module) error

// Module
// Any struct that implements these methods can be considered a module.
type Module interface {
	ParseConfig(map[string]interface{}) error
	Config() ModuleConfig
	GetName() string
	Setup() error
	TearDown() error
	Update() error
}

// ModuleConfig
// General items we may need to track for each module.
type ModuleConfig struct {
	Enabled  bool        `yaml:"enabled"`
	Settings interface{} `yaml:"settings"`
}

// Registry holds Modules and their respective dependencies.
// Modules is a map where keys are module identifiers and values are Module instances.
// Dependencies is a map where keys are module identifiers and values are slices of module identifiers that the key module depends on.
type Registry struct {
	Modules map[string]Module
}

func NewModuleRegistry() *Registry {
	return &Registry{
		Modules: make(map[string]Module),
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
func (mr *Registry) Register(mod Module) {
	mr.Modules[mod.GetName()] = mod
}

// RunSetup Runs the Setup method on each Registry.Modules
func (mr *Registry) RunSetup() (err error) {
	return mr.execute(Module.Setup)
}

// RunUpdate runs the Setup method on each Registry.Modules
func (mr *Registry) RunUpdate() (err error) {
	return mr.execute(Module.Update)
}

// RunTeardown runs the Setup method on each Registry.Modules
func (mr *Registry) RunTeardown() (err error) {
	return mr.execute(Module.TearDown)
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
		rawSettings, ok := moduleSettingsMap[moduleName]
		if !ok {
			return fmt.Errorf("unable to get raw settings for module: %s", moduleName)
		}

		rawSettingsMap, ok := rawSettings.(map[string]interface{})
		if !ok {
			return fmt.Errorf("settings for module %s are not a valid map", moduleName)
		}

		err = module.ParseConfig(rawSettingsMap)
		if err != nil {
			return err
		}
		if rawSettings, ok := moduleSettingsMap[moduleName]; ok {
			err = module.ParseConfig(rawSettings.(map[string]interface{}))
			if err != nil {
				return fmt.Errorf("parseConfig failed parsing: %w", err)
			}
		}
	}

	return nil
}

// execute runs the action on all modules in no order.
func (mr *Registry) execute(action moduleAction) (err error) {
	for _, module := range mr.Modules {
		if !module.Config().Enabled {
			return nil
		}
		err = action(module)
		if err != nil {
			return fmt.Errorf("RunSetup failed: %w", err)
		}
	}
	return
}
