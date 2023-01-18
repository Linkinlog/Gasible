package cmd

import (
	"github.com/Linkinlog/gasible/cmd/yamlParser"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(generatorCmd)
}

var generatorCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a default YAML file.",
	Long: `This will create a default YAML file using the 
  sensible defaults we provided.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := yamlParser.CreateDefaults()
		if err != nil {
			panic(err)
		}
	},
}
