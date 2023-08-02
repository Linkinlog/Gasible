package cmd

import (
	"github.com/Linkinlog/gasible/internal/app"
	"github.com/spf13/cobra"
)

func newUpdateCmd(app *app.App) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "update",
		Short: "Update packages and configurations.",
		Long:  `This will run the update command against all modules.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := app.ModuleRegistry.ReadAndSetRegistryConfigsFromYAML()
			if err != nil {
				return err
			}
			return app.ModuleRegistry.RunUpdate()
		},
	})
}
