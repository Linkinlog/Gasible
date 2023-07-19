package cmd

import (
	"github.com/Linkinlog/gasible/internal/core"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update packages and configurations.",
	Long: `This will run the update command against all
    modules.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := core.CurrentState.ReadConfigFromFile("")
		if err != nil {
			return err
		}
		return core.ModuleRegistry.RunUpdate()
	},
}
