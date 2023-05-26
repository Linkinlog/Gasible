package cmd

import (
	"github.com/Linkinlog/gasible/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update packages and configurations.",
	Long: `This will run the update command against all
    modules.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.ModuleRegistry.RunUpdate()
		if err != nil {
			panic(err)
		}
	},
}
