package cmd

import (
	"github.com/Linkinlog/gasible/internal/app"
	"github.com/spf13/cobra"
)

func newTeardown(app *app.App) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "teardown",
		Short: "Teardown all modules.",
		Long:  `This will run the teardown method on all modules, this can result in data/package loss.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := app.ModuleRegistry.ReadAndSetRegistryConfigsFromYAML()
			if err != nil {
				return err
			}
			return app.ModuleRegistry.RunTeardown()
		},
	})
}
