package console

import (
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"github.com/strider2038/digen/pkg/di"
)

func newGenerator(options *Options, config *viper.Viper) (*di.Generator, error) {
	dir := config.GetString("di.dir")
	if dir == "" {
		dir = config.GetString("work_dir")
	}
	configFile := config.GetString("di.config_file")
	if configFile == "" {
		configFile = strings.TrimRight(dir, "/") + "/internal/_config.go"
	}

	generator := &di.Generator{
		BaseDir:        dir,
		ConfigFilename: configFile,
		Version:        options.Version,
		BuildTime:      options.BuildTime,
		Logger:         terminalLogger{},
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
