package app

import (
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"github.com/strider2038/digen/internal/di"
)

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
