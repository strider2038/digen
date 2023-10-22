package console

import (
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"github.com/strider2038/digen/internal/di"
)

func newGenerator(options *Options, config *viper.Viper) (*di.Generator, error) {
	dir := config.GetString("di.dir")
	if dir == "" {
		dir = config.GetString("work_dir")
	}

	generator := &di.Generator{
		BaseDir:   dir,
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

func (log terminalLogger) Warning(a ...interface{}) {
	pterm.Warning.Println(a...)
}
