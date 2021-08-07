package console

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/mod/modfile"
)

type Configuration struct {
	WorkDir           string
	ModulePath        string
	ContainerFilename string
}

func (c *Configuration) RootPackage() string {
	return c.ModulePath + "/" + c.WorkDir
}

var errMissingModule = errors.New("cannot detect module from go.mod")

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

	mod, err := os.ReadFile("go.mod")
	if err != nil {
		return nil, errors.Wrap(err, "failed to read go.mod file")
	}
	path := modfile.ModulePath(mod)
	if path == "" {
		return nil, errMissingModule
	}

	config.ModulePath = path

	return config, nil
}
