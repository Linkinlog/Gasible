package cmd

import (
	"github.com/Linkinlog/gasible/cmd/commandProcessor"
	"github.com/Linkinlog/gasible/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initializer)
}

var initializer = &cobra.Command{
	Use:   "init",
	Short: "Initialize standard setup",
	Long:  `This will read from the config file and set up the system accordingly.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a config struct and fill it from the config file.
		conf := core.CoreConfig{}
		err := commandProcessor.InitProcess(&conf)
		if err != nil {
			return err
		}
		return nil
	},
}
