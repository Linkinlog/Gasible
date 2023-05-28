package cmd

import (
	"github.com/Linkinlog/gasible/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(config)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(teardown)
	rootCmd.AddCommand(versionCmd)
}

var (
	rootCmd = &cobra.Command{
		Use:   "gas",
		Short: "A lightweight configurator for local development environments",
		Long: `Gasible is a tool that can be used to automate the installation of any tool from your favorite OS/Package manager,
        it also provides tooling for setting up bare git repos that can be useful with local configs. Read more in the README.md of this package`,
	}
)

// Execute executes the root command.
func Execute() error {
	err := core.CurrentState.ReadConfigFromFile("")
	if err != nil {
		panic(err)
	}
	return rootCmd.Execute()
}
