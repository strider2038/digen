package console

import (
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"github.com/strider2038/digen/pkg/di"
)

func newGenerator(options *Options, config *viper.Viper) (*di.Generator, error) {
	generator := &di.Generator{
		WorkDir:   config.GetString("work_dir"),
		Version:   options.Version,
		BuildTime: options.BuildTime,
		Logger:    terminalLogger{},
	}

	return generator, nil
}

type terminalLogger struct{}

func (log terminalLogger) Info(a ...interface{}) {
	pterm.Info.Println(a...)
}

func (log terminalLogger) Success(a ...interface{}) {
	pterm.Success.Println(a...)
}
