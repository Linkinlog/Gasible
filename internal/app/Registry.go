package app

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Registry holds Modules and their respective dependencies.
// Modules is a map where keys are module identifiers and values are Module instances.
// Dependencies is a map where keys are module identifiers and values are slices of module identifiers that the key module depends on.
type Registry struct {
	Modules     map[string]Module
	SettingsMap map[string]interface{}
}

func NewRegistry() *Registry {
	return &Registry{
		Modules:     make(map[string]Module),
		SettingsMap: make(map[string]interface{}),
	}
}

func (r *Registry) Register(mod Module) {
	r.Modules[mod.GetName()] = mod
}

func (r *Registry) GetModule(mod string) (Module, error) {
	found := r.Modules[mod]
	if found != nil {
		return found, nil
	}
	return nil, ModuleNotFoundError
}

// TODO abstract this out to Config.go
func (r *Registry) WriteRegistryConfigsToYAML() error {
	r.updateSettingsMap()

	settingsYAML, err := yaml.Marshal(r.SettingsMap)
	if err != nil {
		return err
	}

	filePath, err := createAndOrGetConfigPath()
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, settingsYAML, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *Registry) ReadRegistryConfigsFromYAML() error {
	filePath, err := createAndOrGetConfigPath()
	if err != nil {
		return fmt.Errorf("ReadRegistryConfigsFromYAML error: %w", err)
	}

	fileContents, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return fmt.Errorf("ReadRegistryConfigsFromYAML error: %w", readErr)
	}

	unmarshalErr := yaml.Unmarshal(fileContents, &r.SettingsMap)
	if unmarshalErr != nil {
		return fmt.Errorf("ReadRegistryConfigsFromYAML error: %w", unmarshalErr)
	}

	return nil
}

func (r *Registry) ReadAndSetRegistryConfigsFromYAML() error {
	readErr := r.ReadRegistryConfigsFromYAML()
	if readErr != nil {
		return fmt.Errorf("ReadAndSetRegistryConfigsFromYAML error: %w", readErr)
	}
	for moduleName, module := range r.Modules {
		setErr := r.setCurrent(moduleName, module)
		if setErr != nil {
			return fmt.Errorf("ReadAndSetRegistryConfigsFromYAML error: %w", setErr)
		}
	}
	return nil
}

func (r *Registry) setCurrent(moduleName string, module Module) error {
	rawSettings, ok := r.SettingsMap[moduleName]
	if !ok {
		return ErrorAs("setCurrent", fmt.Errorf("unable to get raw settings for module: %s", moduleName))
	}

	rawSettingsMap, ok := rawSettings.(map[string]interface{})
	if !ok {
		return ErrorAs("setCurrent", fmt.Errorf("settings for module %s are not a valid map", moduleName))
	}

	parseErr := module.ParseConfig(rawSettingsMap)
	if parseErr != nil {
		return parseErr
	}

	return nil
}

func (r *Registry) RunSetup() (err error) {
	return r.execute(Module.Setup)
}

func (r *Registry) RunUpdate() (err error) {
	return r.execute(Module.Update)
}

func (r *Registry) RunTeardown() (err error) {
	return r.execute(Module.TearDown)
}

func (r *Registry) updateSettingsMap() {
	for moduleName, mod := range r.Modules {
		r.SettingsMap[moduleName] = mod.Config()
	}
}

type moduleAction func(Module) error

func (r *Registry) execute(action moduleAction) (err error) {
	for _, module := range r.Modules {
		if !module.Config().Enabled {
			return nil
		}
		err = action(module)
		if err != nil {
			return ErrorAs("execute", err)
		}
	}
	return
}
