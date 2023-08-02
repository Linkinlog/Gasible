package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newVersionCmd() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of Gasible.",
		Long:  `All software has versions. This is ours.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version: v0.1.3")
		},
	})
}
