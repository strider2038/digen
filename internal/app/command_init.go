package app

import (
	"github.com/muonsoft/errors"
	"github.com/strider2038/digen/internal/config"
)

func runInit(options *Options) error {
	params, err := config.Init()
	if err != nil {
		return errors.Errorf("init config: %w", err)
	}

	return newGenerator(options, params).Initialize()
}
