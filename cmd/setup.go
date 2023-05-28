package cmd

import (
	"github.com/Linkinlog/gasible/internal/core"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	RootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Run the setup for all modules.",
	Run: func(cmd *cobra.Command, args []string) {
		err := core.CurrentState.Config.ModuleRegistry.RunSetup()
		if err != nil {
			log.Fatal(err)
		}
	},
}
