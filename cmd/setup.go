package cmd

import (
	"github.com/Linkinlog/gasible/internal/core"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up all modules.",
	Long:  `This will run the setup method on all modules.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := core.CurrentState.ReadConfigFromFile("")
		if err != nil {
			return err
		}
		return core.CurrentState.Config.ModuleRegistry.RunSetup()
	},
}
