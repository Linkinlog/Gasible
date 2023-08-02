package cmd

import (
	"github.com/Linkinlog/gasible/internal/app"
	"github.com/spf13/cobra"
)

func newWriteCurrent(app *app.App) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "generate",
		Short: "Write the current config to $HOME/.gas/config.yml.",
		Long: `This will create a default YAML file using the 
  defaults provided by each module.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.ModuleRegistry.WriteRegistryConfigsToYAML()
		},
	})
}
