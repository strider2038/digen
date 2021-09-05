package console

import (
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"github.com/strider2038/digen"
)

func newGenerator(options *Options) (*digen.Generator, error) {
	v := viper.New()
	v.SetConfigName("digen")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	config := &digen.Generator{
		WorkDir:   v.GetString("work_dir"),
		Version:   options.Version,
		BuildTime: options.BuildTime,
		Logger:    terminalLogger{},
	}

	return config, nil
}

type terminalLogger struct{}

func (log terminalLogger) Info(a ...interface{}) {
	pterm.Info.Println(a...)
}

func (log terminalLogger) Success(a ...interface{}) {
	pterm.Success.Println(a...)
}
