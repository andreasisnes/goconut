package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove usage",
	Short: "remove",
	Long:  `remove long`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.PrintErrln("Invalid number of arguments")
			os.Exit(1)
		}

		parseModuleName(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		secrets, err := readSecrets()
		if err != nil {
			exit(cmd, err.Error(), 1)
		}

		delete(secrets, key)
		if _, err := dumpSecrets(secrets); err != nil {
			exit(cmd, err.Error(), 1)
		}
	},
}

func init() {
}
