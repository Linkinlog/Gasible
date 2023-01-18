package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	CfgFile string
	Noop    bool

	RootCmd = &cobra.Command{
		Use:   "gasible",
		Short: "A lightweight configurator for local development environments",
		Long: `Gasible is a tool that can be used to automate the installation of any tool from your favorite OS/Package manager,
        it also provides tooling for setting up bare git repos that can be useful with local configs. Read more in the README.md of this package`,
	}
)

// Execute executes the root command.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.PersistentFlags().StringVar(&CfgFile, "config", ".gas.yml", "config file")
	RootCmd.PersistentFlags().BoolVar(&Noop, "noop", false, "No-Op command to run command without making any changes")
}
