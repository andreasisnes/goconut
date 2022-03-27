package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear usage",
	Short: "clear",
	Long:  `clear long`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			cmd.PrintErrln("Invalid number of arguments")
			os.Exit(1)
		}

		parseModuleName(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.Remove(getModuleSecretspath()); err != nil {
			cmd.PrintErrln("Invalid number of arguments")
			os.Exit(1)
		}
	},
}

func init() {
}
