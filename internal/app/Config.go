package app

import (
	"os"
	"path/filepath"
)

type Config struct {
	Version    string   `yaml:"version"`
	AllModules []Module `yaml:"modules"`
	FullPath   string
	// TODO log level/filepath once logging is implemented
}

func NewConfig() *Config {
	return &Config{
		Version:    "0.1.0",
		AllModules: make([]Module, 0),
		FullPath:   mustGetConfigPath(),
	}
}

// createAndOrGetConfigPath
// Creates the structure for the config if needed.
// Returns full os compliant path once found.
func createAndOrGetConfigPath() (string, error) {
	// Find home directory.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", ErrorAs("createAndOrGetConfigPath", err)
	}

	// Append the AppConfig file directory to the home directory path
	confDir := filepath.Join(homeDir, configDir)
	_, err = os.Stat(confDir)

	// If the directory doesn't exist, we create it
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(confDir, 0755)
		if errDir != nil {
			return "", ErrorAs("createAndOrGetConfigPath", err)
		}
	}

	// Append the AppConfig file path to the home directory path
	confFilePath := filepath.Join(confDir, configFilename)

	_, err = os.Stat(confFilePath)

	// If the config doesn't exist, we create it
	if os.IsNotExist(err) {
		_, errDir := os.Create(confFilePath)
		if errDir != nil {
			return "", ErrorAs("createAndOrGetConfigPath", err)
		}
	}

	return confFilePath, nil
}

func mustGetConfigPath() string {
	path, err := createAndOrGetConfigPath()
	if err != nil {
		panic(err)
	}
	return path
}
