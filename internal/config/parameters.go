package config

import (
	"io/fs"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/muonsoft/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"github.com/strider2038/digen/internal/di"
)

var errInvalidPath = errors.New("invalid path")

type Parameters struct {
	Version       string        `json:"version" yaml:"version"`
	Container     Container     `json:"container" yaml:"container"`
	ErrorHandling ErrorHandling `json:"errorHandling" yaml:"errorHandling"`
}

type Container struct {
	Dir string `json:"dir" yaml:"dir"`
}

type ErrorHandling struct {
	New  ErrorOptions     `json:"new" yaml:"new"`
	Join ErrorOptions     `json:"join" yaml:"join"`
	Wrap WrapErrorOptions `json:"wrap" yaml:"wrap"`
}

func (h ErrorHandling) MapToOptions() di.ErrorHandling {
	return di.ErrorHandling{
		New:  h.New.mapToOptions(),
		Join: h.Join.mapToOptions(),
		Wrap: h.Wrap.mapToOptions(),
	}
}

type ErrorOptions struct {
	Pkg  string `json:"pkg" yaml:"pkg"`
	Func string `json:"func" yaml:"func"`
}

func (o ErrorOptions) mapToOptions() di.ErrorOptions {
	return di.ErrorOptions{
		Package:  o.Pkg,
		Function: o.Func,
	}
}

type WrapErrorOptions struct {
	Pkg  string `json:"pkg" yaml:"pkg"`
	Func string `json:"func" yaml:"func"`
	Verb string `json:"verb" yaml:"verb"`
}

func (o WrapErrorOptions) mapToOptions() di.ErrorOptions {
	return di.ErrorOptions{
		Package:  o.Pkg,
		Function: o.Func,
		Verb:     o.Verb,
	}
}

func Load() (*Parameters, error) {
	config := newConfig()
	err := config.ReadInConfig()
	if err != nil {
		return nil, errors.Errorf("read config: %w", err)
	}

	return readParameters(config)
}

func Init() (*Parameters, error) {
	config := newConfig()
	err := config.ReadInConfig()
	if err != nil {
		if errors.IsOfType[viper.ConfigFileNotFoundError](err) {
			if err := initConfig(config); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.Errorf("read config: %w", err)
		}
	}

	return readParameters(config)
}

func readParameters(config *viper.Viper) (*Parameters, error) {
	var parameters Parameters
	if err := config.Unmarshal(&parameters); err != nil {
		return nil, errors.Errorf("unmarshal config: %w", err)
	}
	if err := checkVersion(parameters.Version); err != nil {
		return nil, err
	}

	return &parameters, nil
}

func newConfig() *viper.Viper {
	config := viper.New()
	config.SetConfigName("digen")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")

	return config
}

func initConfig(config *viper.Viper) error {
	if isFileExist("digen.yaml") {

		return nil
	}

	prompt := promptui.Prompt{
		Label:   "enter path to working directory",
		Default: "di",
		Validate: func(path string) error {
			if fs.ValidPath(path) {
				return nil
			}

			return errInvalidPath
		},
	}

	dir, err := prompt.Run()
	if err != nil {
		return err
	}

	config.Set("version", Version)
	config.Set("container.dir", dir)
	config.Set("errorHandling.package", "errors")
	config.Set("errorHandling.wrapPackage", "fmt")
	config.Set("errorHandling.wrapFunction", "Errorf")
	err = config.SafeWriteConfig()
	if err != nil {
		return errors.Errorf("write config: %w", err)
	}

	pterm.Success.Println("configuration file generated: digen.yaml")

	return nil
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}
