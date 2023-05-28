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
		return core.CurrentState.Config.ModuleRegistry.RunSetup()
	},
}
