package config

import (
	"io/fs"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/muonsoft/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

type Parameters struct {
	AppVersion string     `json:"app_version" yaml:"app_version"`
	Containers Containers `json:"containers" yaml:"containers"`
}

type Containers struct {
	Dir string `json:"dir" yaml:"dir"`
}

func Load() (*Parameters, error) {
	config := newConfig()
	err := config.ReadInConfig()
	if err != nil {
		return nil, errors.Errorf("read config: %w", err)
	}
}

func Init() (*Parameters, error) {
	config := newConfig()
	err := config.ReadInConfig()
	if err != nil {
		if errors.IsOfType[viper.ConfigFileNotFoundError](err) {
			if err := initConfig(config); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.Errorf("read config: %w", err)
		}
	}
}

const (
	configAppVersion = "app_version"
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
