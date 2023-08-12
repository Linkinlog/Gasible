package app

// configDir is so that we can specify where to put our config.
const configDir = ".gas"

// configFileName is so that we can specify which filename our config should use.
const configFilename = "config.yml"

// App is to represent the currently running application and its state.
type App struct {
	Config         *Config
	ModuleRegistry *registry
	Version        string
}

// New returns a pointer to an application
func New() *App {
	return &App{
		Config:         NewConfig(),
		ModuleRegistry: newRegistry(),
		Version:        "0.1.3",
	}
}
