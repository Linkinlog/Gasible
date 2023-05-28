package core

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// CurrentState
// The state of the app.
var CurrentState = coreConfig{
	Config: Config{ModuleRegistry: ModuleRegistry},
}

// moduleSettingsMap is a map holding settings for each module.
var moduleSettingsMap = make(map[string]interface{})

const configDir = ".gas"
const configFile = ".config.yml"

// coreConfig
// The overall config for the app.
type coreConfig struct {
	Config
}

// Config The entire Config YAML.
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Config struct {
	ModuleRegistry *Registry
}

// WriteConfigToFile will generate a YAML file
// using the defaults we outline.
func (conf *Config) WriteConfigToFile() (err error) {
	// For each module in the registry,
	// retrieve its settings and store them in the ModuleSettings map.
	for moduleName, module := range conf.ModuleRegistry.Modules {
		moduleSettingsMap[moduleName] = module.Config().Settings
	}

	// Marshal the ModuleSettings map to YAML.
	settingsYAML, err := yaml.Marshal(moduleSettingsMap)
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

// ReadConfigFromFile
// Attempts to read the file at filePath and set the settings for each module
// using the file and yaml.Unmarshal.
func (conf *Config) ReadConfigFromFile(filePath string) (err error) {
	var fileContents []byte
	if filePath == "" {
		filePath, err = getConfigPath()
	}
	if err != nil {
		return err
	}

	fileContents, err = os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return conf.ModuleRegistry.setCurrent(fileContents)
}

// getConfigPath
// Creates the structure for the config if needed.
// Returns full os compliant path once found.
func getConfigPath() (string, error) {
	// Find home directory.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Append the Config file directory to the home directory path
	confDir := filepath.Join(homeDir, configDir)
	_, err = os.Stat(confDir)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(confDir, 0755)
		if errDir != nil {
			return "", err
		}
	}

	// Append the Config file path to the home directory path
	confFile := filepath.Join(confDir, configFile)

	_, err = os.Stat(confFile)

	if os.IsNotExist(err) {
		_, errDir := os.Create(confFile)
		if errDir != nil {
			return "", err
		}
	}

	return confFile, nil
}
