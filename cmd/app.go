package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	moduleFlag    = "module"
	moduleNameKey = "moduleName"
	projectMapDir = "usersecrets"
)

var App = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  ``,
	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
	},
	Run: func(ccmd *cobra.Command, args []string) {
		ccmd.HelpFunc()(ccmd, args)
	},
}

func init() {
	App.PersistentFlags().StringP("module", "m", "", "Go module file")
	App.AddCommand(setCmd)
	App.AddCommand(clearCmd)
	App.AddCommand(listCmd)
	App.AddCommand(removeCmd)
}

func exit(cmd *cobra.Command, message string, statusCode int) {
	cmd.PrintErrln(message)
	os.Exit(statusCode)
}
