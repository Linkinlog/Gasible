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
	return mr.executeInOrder(Module.Setup)
}

// RunUpdate runs the Setup method on each Registry.Modules
func (mr *Registry) RunUpdate() (err error) {
	return mr.executeInOrder(Module.Update)
}

// RunTeardown runs the Setup method on each Registry.Modules
func (mr *Registry) RunTeardown() (err error) {
	return mr.executeInOrder(Module.TearDown)
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

// executeInOrder runs the action on all modules in topologically sorted order.
func (mr *Registry) executeInOrder(action moduleAction) (err error) {
	order, err := mr.TopologicallySortedModuleDeps()
	if err != nil {
		return fmt.Errorf("could not order: %w", err)
	}
	for _, moduleName := range order {
		module := mr.Modules[moduleName]
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

// TopologicallySortedModuleDeps sorts modules into a directed acyclic graph (hopefully).
func (mr *Registry) TopologicallySortedModuleDeps() ([]string, error) {
	order := make([]string, 0)
	visited := make(map[string]bool)
	temp := make(map[string]bool)

	var visitAllDependencies func(node string) error
	visitAllDependencies = func(node string) error {
		if temp[node] {
			return fmt.Errorf("dependency %s exists, %w", node, dependencyGraphCycleError)
		}
		if !visited[node] {
			temp[node] = true
			for _, v := range mr.Dependencies[node] {
				if _, exists := mr.Dependencies[v]; !exists {
					return fmt.Errorf("dependency %s for module %s does not exist", v, node)
				}
				err := visitAllDependencies(v)
				if err != nil {
					return err
				}
			}
			visited[node] = true
			temp[node] = false
			order = append(order, node)
		}
		return nil
	}

	for k := range mr.Dependencies {
		err := visitAllDependencies(k)
		if err != nil {
			return nil, err
		}
	}

	return order, nil
}
