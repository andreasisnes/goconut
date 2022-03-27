package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	Long:  `set value`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			cmd.PrintErrln("Invalid number of arguments")
			os.Exit(1)
		}

		parseModuleName(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		if err := initializeTree(getModuleDir(), getModuleSecretspath(), make(map[string]interface{})); err != nil {
			exit(cmd, err.Error(), 1)
		}

		secrets, err := readSecrets()
		if err != nil {
			exit(cmd, err.Error(), 1)
		}

		secrets[key] = value
		if _, err := dumpSecrets(secrets); err != nil {
			exit(cmd, err.Error(), 1)
		}
	},
}

func init() {
}
