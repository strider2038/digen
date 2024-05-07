package app

import (
	"github.com/muonsoft/errors"
	"github.com/spf13/viper"
)

func runInit(options *Options) error {

	return newGenerator(options, config).Initialize()
}
