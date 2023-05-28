package cmd

import (
	"github.com/Linkinlog/gasible/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(config)
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
  sensible defaults we provided.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.WriteConfigToFile(&core.CurrentState.Config)
		if err != nil {
			panic(err)
		}
	},
}
