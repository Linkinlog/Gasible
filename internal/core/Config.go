package core

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var CurrentState = CoreConfig{
	Config: Config{ModuleRegistry: ModuleRegistry},
}

// moduleSettings is a map holding settings for each module.
var moduleSettings = make(map[string]interface{})

// CoreConfig The state of the app.
type CoreConfig struct {
	Config
}

// Config The entire config YAML.
type Config struct {
	ModuleRegistry *Registry
}

// WriteConfigToFile will generate a YAML file
// using the defaults we outline.
func WriteConfigToFile(model *Config) (err error) {
	// For each module in the registry,
	// retrieve its settings and store them in the ModuleSettings map.
	for moduleName, module := range CurrentState.Config.ModuleRegistry.Modules {
		moduleSettings[moduleName] = module.Config().Settings
	}

	// Marshal the ModuleSettings map to YAML.
	settingsYAML, err := yaml.Marshal(moduleSettings)
	if err != nil {
		return err
	}

	filePath, err := getConfigPath()
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, settingsYAML, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadConfigFromFile(filePath string) (err error) {
	var fileContents []byte
	if filePath == "" {
		filePath, err = getConfigPath()
	}

	fileContents, err = os.ReadFile(filePath)

	return CurrentState.Config.ModuleRegistry.setCurrent(fileContents)
}

func getConfigPath() (string, error) {
	// Find home directory.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	confDir := filepath.Join(homeDir, ".gas")
	// Append the config file path to the home directory path
	_, err = os.Stat(confDir)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(confDir, 0755)
		if errDir != nil {
			return "", err
		}
	}
	return filepath.Join(homeDir, ".gas/config.yml"), nil
}
