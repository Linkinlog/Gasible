package cmd

import (
	"github.com/Linkinlog/gasible/cmd/commandProcessor"
	"github.com/Linkinlog/gasible/internal/models"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initializer)
}

var initializer = &cobra.Command{
	Use:   "init",
	Short: "Initialize standard setup",
	Long:  `This will read from the config file and set up the system accordingly.`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("config")
		noop, _ := cmd.Flags().GetBool("noop")
		// Create a config struct and fill it from the config file.
		conf := models.Config{}
		// Overwrite GlobalOpts with our defaults
		conf.GlobalOpts.FilePath = filePath
		conf.FillFromFile()
		conf.GlobalOpts.NoOp = noop
		err := commandProcessor.InitProcess(&conf)
		if err != nil {
			panic(err)
		}
	},
}
