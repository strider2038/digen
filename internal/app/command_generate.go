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
		BaseDir:   params.Container.Dir,
		Version:   options.Version,
		BuildTime: options.BuildTime,
		Logger:    terminalLogger{},
		ErrorWrapping: di.ErrorHandling{
			Package:      params.ErrorHandling.Package,
			WrapPackage:  params.ErrorHandling.WrapPackage,
			WrapFunction: params.ErrorHandling.WrapFunction,
			Verb:         params.ErrorHandling.Verb,
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
