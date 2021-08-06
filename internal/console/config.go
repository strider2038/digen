package console

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Configuration struct {
	WorkDir           string
	ContainerFilename string
}

func loadConfig() (*Configuration, error) {
	v := viper.New()
	v.SetConfigName("digen")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	config := &Configuration{
		WorkDir: v.GetString("work_dir"),
	}

	if config.WorkDir == "" {
		config.WorkDir = "."
	}
	if config.ContainerFilename == "" {
		config.ContainerFilename = config.WorkDir + "/internal/container.go"
	}

	return config, nil
}
