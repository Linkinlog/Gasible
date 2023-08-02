package app

const configDir = ".gas"
const configFilename = "config.yml"

type App struct {
	Config         *Config
	ModuleRegistry *registry
}

// New returns a pointer to an application
func New() *App {
	return &App{
		Config:         NewConfig(),
		ModuleRegistry: newRegistry(),
	}
}
