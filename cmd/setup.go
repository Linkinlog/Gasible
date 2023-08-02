package cmd

import (
	"github.com/Linkinlog/gasible/internal/app"
	"github.com/spf13/cobra"
)

func newSetupCmd(app *app.App) {
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "setup",
			Short: "Set up all modules.",
			Long:  `This will run the setup method on all modules.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				err := app.ModuleRegistry.ReadAndSetRegistryConfigsFromYAML()
				if err != nil {
					return err
				}
				return app.ModuleRegistry.RunSetup()
			},
		},
	)
}
