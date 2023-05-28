package cmd

import (
	"github.com/Linkinlog/gasible/internal/core"
	"github.com/spf13/cobra"
)

var teardown = &cobra.Command{
	Use:   "teardown",
	Short: "Teardown all modules.",
	Long:  `This will run the teardown method on all modules, this can result in data/package loss.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.CurrentState.Config.ModuleRegistry.RunTeardown()
	},
}
