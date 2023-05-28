package cmd

import (
	"github.com/Linkinlog/gasible/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	config.AddCommand(writeCurrent)
}

var config = &cobra.Command{
	Use:   "config",
	Short: "Manage the config.yml",
}

var writeCurrent = &cobra.Command{
	Use:   "generate",
	Short: "Write the current config to ./config.yml.",
	Long: `This will create a default YAML file using the 
  defaults provided by each module.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.CurrentState.Config.WriteConfigToFile()
	},
}
