package app

import "errors"

var ModuleNotFoundError = errors.New("no module found")

// Module
// Any struct that implements these methods can be considered a module.
type Module interface {
	ParseConfig(map[string]interface{}) error
	Config() ModuleConfig
	GetName() string
	Setup() error
	TearDown() error
	Update() error
	SetApp(app *App)
}

// ModuleConfig
// General items we may need to track for each module.
type ModuleConfig struct {
	Enabled  bool        `yaml:"enabled"`
	Settings interface{} `yaml:"settings"`
}
