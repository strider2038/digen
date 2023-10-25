package app

import (
	"io/fs"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/muonsoft/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

const (
	configAppVersion   = "app_version"
	configContainerDir = "container.dir"
)

var errInvalidPath = errors.New("invalid path")

func newConfig() *viper.Viper {
	config := viper.New()
	config.SetConfigName("digen")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")

	return config
}

func initConfig(config *viper.Viper) error {
	if isFileExist("digen.yaml") {

		return nil
	}

	prompt := promptui.Prompt{
		Label:   "enter path to working directory",
		Default: "di",
		Validate: func(path string) error {
			if fs.ValidPath(path) {
				return nil
			}

			return errInvalidPath
		},
	}

	dir, err := prompt.Run()
	if err != nil {
		return err
	}

	config.Set(configAppVersion, Version)
	config.Set(configContainerDir, dir)
	err = config.SafeWriteConfig()
	if err != nil {
		return errors.Errorf("write config: %w", err)
	}

	pterm.Success.Println("configuration file generated: digen.yaml")

	return nil
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}
