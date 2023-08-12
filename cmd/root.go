package cmd

import (
	"github.com/Linkinlog/gasible/internal/app"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gasible",
		Short: "A lightweight configurator for local development environments",
		Long: `Gasible is a tool that can be used to automate the installation of any tool from your favorite
OS/Package manager, it also provides tooling for setting up bare git repos that can be useful with
local configs. Read more in the README.md of this package`,
	}
)

func ExecuteApplication(app *app.App) error {
	newVersionCmd(app)
	newWriteCurrent(app)
	newSetupCmd(app)
	newUpdateCmd(app)
	newTeardown(app)
	return rootCmd.Execute()
}
