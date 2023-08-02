package app

import (
	"fmt"
	"github.com/Linkinlog/gasible/internal"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

// registry holds Modules and their respective dependencies.
// Modules is a map where keys are module identifiers and values are Module instances.
// Dependencies is a map where keys are module identifiers and values are slices of module identifiers that the key module depends on.
type registry struct {
	Modules     map[string]Module
	SettingsMap map[string]interface{}
}

func newRegistry() *registry {
	return &registry{
		Modules:     make(map[string]Module),
		SettingsMap: make(map[string]interface{}),
	}
}

func (r *registry) Register(mod Module) {
	r.Modules[mod.GetName()] = mod
}

func (r *registry) GetModule(mod string) Module {
	return r.Modules[mod]
}

// TODO abstract write/read out to Config.go

func (r *registry) WriteRegistryConfigsToYAML() error {
	r.updateSettingsMap()

	settingsYAML, err := yaml.Marshal(r.SettingsMap)
	if err != nil {
		return err
	}

	filePath, err := createAndOrGetConfigPath()
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, settingsYAML, 0600)
	if err != nil {
		return err
	}

	log.Printf("Config successfully generated to %s, have fun!\n", filePath)
	return nil
}

func (r *registry) readRegistryConfigsFromYAML() error {
	filePath, err := createAndOrGetConfigPath()
	if err != nil {
		return fmt.Errorf("ReadRegistryConfigsFromYAML error: %w", err)
	}

	fileContents, readErr := os.ReadFile(filepath.Clean(filePath))
	if readErr != nil {
		return fmt.Errorf("ReadRegistryConfigsFromYAML error: %w", readErr)
	}

	unmarshalErr := yaml.Unmarshal(fileContents, &r.SettingsMap)
	if unmarshalErr != nil {
		return fmt.Errorf("ReadRegistryConfigsFromYAML error: %w", unmarshalErr)
	}

	return nil
}

func (r *registry) ReadAndSetRegistryConfigsFromYAML() error {
	readErr := r.readRegistryConfigsFromYAML()
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

func (r *registry) setCurrent(moduleName string, module Module) error {
	var settingsNotFoundErr error = fmt.Errorf("unable to get raw settings for module: %s", moduleName)
	var settingsNotValidErr error = fmt.Errorf("settings for module %s are not a valid map", moduleName)
	rawSettings, ok := r.SettingsMap[moduleName]
	if !ok {
		return internal.ErrorAs("setCurrent", settingsNotFoundErr)
	}

	rawSettingsMap, ok := rawSettings.(map[string]interface{})
	if !ok {
		return internal.ErrorAs("setCurrent", settingsNotValidErr)
	}

	parseErr := module.ParseConfig(rawSettingsMap)
	if parseErr != nil {
		return parseErr
	}

	return nil
}

func (r *registry) RunSetup() (err error) {
	return r.execute(Module.Setup)
}

func (r *registry) RunUpdate() (err error) {
	return r.execute(Module.Update)
}

func (r *registry) RunTeardown() (err error) {
	return r.execute(Module.TearDown)
}

func (r *registry) updateSettingsMap() {
	for moduleName, mod := range r.Modules {
		r.SettingsMap[moduleName] = mod.Config()
	}
}

type moduleAction func(Module) error

func (r *registry) execute(action moduleAction) (err error) {
	for _, module := range r.Modules {
		if !module.Config().Enabled {
			return nil
		}
		err = action(module)
		if err != nil {
			return internal.ErrorAs("registry.execute", err)
		}
	}
	return
}
