package cmd

import (
	"encoding/json"
	"errors"
	"go/build"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/mod/modfile"
)

func getBaseDirpath() string {
	return path.Join(build.Default.GOPATH, projectMapDir)
}

func getModulename() string {
	return viper.GetString(moduleNameKey)
}

func getModuleDir() string {
	return path.Join(getBaseDirpath(), getModulename())
}

func getModuleSecretspath() string {
	return path.Join(getModuleDir(), "secrets.json")
}

func readSecrets() (map[string]interface{}, error) {
	content, err := ioutil.ReadFile(getModuleSecretspath())
	if err != nil {
		return nil, err
	}

	values := make(map[string]interface{})
	err = json.Unmarshal(content, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func dumpSecrets(secrets map[string]interface{}) (map[string]interface{}, error) {
	content, err := json.Marshal(secrets)
	if err != nil {
		return nil, err
	}

	return secrets, ioutil.WriteFile(getModuleSecretspath(), content, os.ModePerm)
}

func initializeTree(baseDir, valueFile string, value interface{}) error {
	if _, err := os.Stat(baseDir); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(baseDir, os.ModePerm)
	} else if err != nil {
		return err
	}

	if _, err := os.Stat(valueFile); errors.Is(err, os.ErrNotExist) {
		file, err := os.OpenFile(valueFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			return err
		}

		value, err := json.Marshal(value)
		if err != nil {
			return err
		}

		file.Write(value)
	} else if err != nil {
		return err
	}

	return nil
}

func parseModuleName(cmd *cobra.Command, args []string) {
	if cmd.Flag(moduleFlag).Value.String() == "" {
		if pwd, err := os.Getwd(); err == nil {
			cmd.Flag(moduleFlag).Value.Set(path.Join(pwd, "go.mod"))
		} else {
			exit(cmd, "Failed with: "+err.Error(), 1)
		}
	}

	if stat, err := os.Stat(cmd.Flag(moduleFlag).Value.String()); err != nil {
		exit(cmd, "Failed with: "+err.Error(), 1)
	} else {
		if stat.IsDir() {
			cmd.Flag(moduleFlag).Value.Set(path.Join(cmd.Flag(moduleFlag).Value.String(), "go.mod"))
		}
	}

	if _, err := os.Stat(cmd.Flag(moduleFlag).Value.String()); err == nil {
		if content, err := ioutil.ReadFile(cmd.Flag(moduleFlag).Value.String()); err == nil {
			viper.Set(moduleNameKey, modfile.ModulePath(content))
		} else {
			exit(cmd, "Failed with: "+err.Error(), 1)
		}
	} else {
		exit(cmd, "Failed with: "+err.Error(), 1)
	}
}
