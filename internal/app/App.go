package app

const configDir = ".gas"
const configFilename = "config.yml"

type App struct {
	Config         *Config
	Executor       cmdExecutor
	ModuleRegistry *Registry
	System         *CurrentSystem
}

func New() *App {
	return &App{
		Config:         NewConfig(),
		System:         NewSystem(),
		ModuleRegistry: NewRegistry(),
		Executor:       &NormalRunner{},
	}
}
