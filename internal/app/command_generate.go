package app

import (
	"github.com/muonsoft/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"github.com/strider2038/digen/internal/di"
)

func runGenerate(options *Options) error {

	if err := checkVersion(config.GetString(configAppVersion), Version); err != nil {
		return err
	}

	return newGenerator(options, config).Generate()
}

func newGenerator(options *Options, config *viper.Viper) *di.Generator {
	return &di.Generator{
		BaseDir:   config.GetString(configContainerDir),
		Version:   options.Version,
		BuildTime: options.BuildTime,
		Logger:    terminalLogger{},
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
