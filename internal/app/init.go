package app

import (
	"github.com/muonsoft/errors"
	"github.com/spf13/viper"
)

func runInit(options *Options) error {
	config := newConfig()
	err := config.ReadInConfig()
	if err != nil {
		if errors.IsOfType[viper.ConfigFileNotFoundError](err) {
			if err := initConfig(config); err != nil {
				return err
			}
		} else {
			return errors.Errorf("read config: %w", err)
		}
	}

	return newGenerator(options, config).Initialize()
}
