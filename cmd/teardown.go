package cmd

import (
	"github.com/Linkinlog/gasible/internal/app"
	"github.com/Linkinlog/gasible/internal/modules"
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
				panic(err)
			}
			gpm := app.ModuleRegistry.GetModule("GenericPackageManager").(*modules.GenericPackageManager)
			call := app.ModuleRegistry.GetModule("SysCall").(*modules.SysCall)
			teardownErr := gpm.Manager().Uninstall(modules.ToBeInstalled[gpm.Manager()], call)
			if teardownErr != nil {
				return teardownErr
			}
			return app.ModuleRegistry.RunTeardown()
		},
	})
}
