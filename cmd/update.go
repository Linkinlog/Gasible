package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update packages and configurations.",
	Long: `This will run the update command against all
    packages in the package list, as well as sync the config repo.`,
	Run: func(cmd *cobra.Command, args []string) {
        // TODO
	},
}

