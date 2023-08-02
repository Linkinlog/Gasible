package cmd

import (
	"github.com/Linkinlog/gasible/internal/app"
	"github.com/Linkinlog/gasible/internal/modules"
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
					panic(err)
				}
				gpm := app.ModuleRegistry.GetModule("GenericPackageManager").(*modules.GenericPackageManager)
				call := app.ModuleRegistry.GetModule("SysCall").(*modules.SysCall)
				installErr := gpm.Manager().Install(modules.ToBeInstalled[gpm.Manager()], call)
				if installErr != nil {
					return installErr
				}
				return app.ModuleRegistry.RunSetup()
			},
		},
	)
}
