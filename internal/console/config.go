package console

import (
	"io/fs"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
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
		Label: "Enter path to work directory",
		Validate: func(path string) error {
			if fs.ValidPath(path) || fs.ValidPath(strings.TrimLeft(path, "./")) {
				return nil
			}

			return errInvalidPath
		},
	}

	workDir, err := prompt.Run()
	if err != nil {
		return err
	}

	config.Set("work_dir", workDir)
	err = config.SafeWriteConfig()
	if err != nil {
		return errors.Wrap(err, "failed to write config")
	}

	pterm.Success.Println("configuration file generated: digen.yaml")

	return nil
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}
