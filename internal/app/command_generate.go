package app

import (
	"github.com/muonsoft/errors"
	"github.com/pterm/pterm"
	"github.com/strider2038/digen/internal/config"
	"github.com/strider2038/digen/internal/di"
)

func runGenerate(options *Options) error {
	params, err := config.Load()
	if err != nil {
		return errors.Errorf("load config: %w", err)
	}

	return newGenerator(options, params).Generate()
}

func newGenerator(options *Options, params *config.Parameters) *di.Generator {
	return &di.Generator{
		BaseDir: params.Container.Dir,
		Logger:  terminalLogger{},
		Params: di.GenerationParameters{
			Version:       options.Version,
			ErrorHandling: params.ErrorHandling.MapToOptions(),
		},
	}
}

type terminalLogger struct{}

func (log terminalLogger) Info(a ...interface{}) {
	pterm.Info.Println(a...)
}

func (log terminalLogger) Success(a ...interface{}) {
	pterm.Success.Println(a...)
}

func (log terminalLogger) Warning(a ...interface{}) {
	pterm.Warning.Println(a...)
}
