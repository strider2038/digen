package config

import (
	"encoding/json"
	"io/fs"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/muonsoft/errors"
	"github.com/pterm/pterm"
	"gopkg.in/yaml.v3"
)

var errNoConfig = errors.New("no config file found")

func Load() (*Parameters, error) {
	params, err := loadConfig()
	if err != nil {
		return nil, errors.Errorf("load config: %w", err)
	}

	return params, nil
}

func Init() (*Parameters, error) {
	params, err := loadConfig()
	if errors.Is(err, errNoConfig) {
		params, err = initDefaultConfig()
	}
	if err != nil {
		return nil, errors.Errorf("load config: %w", err)
	}

	return params, nil
}

func initDefaultConfig() (*Parameters, error) {
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
		return nil, err
	}

	params := Parameters{
		Version:   Version,
		Container: Container{Dir: dir},
	}
	data, err := yaml.Marshal(params)
	if err != nil {
		return nil, errors.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile("digen.yaml", data, 0644); err != nil {
		return nil, errors.Errorf("write config: %w", err)
	}

	pterm.Success.Println("configuration file generated: digen.yaml")

	return &params, nil
}

func loadConfig() (*Parameters, error) {
	var (
		content []byte
		err     error
		isYAML  bool
	)
	if isFileExist("digen.yaml") {
		isYAML = true
		content, err = os.ReadFile("digen.yaml")
	} else if isFileExist("digen.yml") {
		isYAML = true
		content, err = os.ReadFile("digen.yml")
	} else if isFileExist("digen.json") {
		content, err = os.ReadFile("digen.json")
	} else {
		return nil, errors.Wrap(errNoConfig)
	}
	if err != nil {
		return nil, errors.Errorf("read config: %w", err)
	}

	var params Parameters
	if isYAML {
		if err := yaml.Unmarshal(content, &params); err != nil {
			return nil, errors.Errorf("unmarshal yaml config: %w", err)
		}
	} else {
		if err := json.Unmarshal(content, &params); err != nil {
			return nil, errors.Errorf("unmarshal json config: %w", err)
		}
	}

	return &params, nil
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}

func isAnyFileExists(filenames ...string) bool {
	for _, filename := range filenames {
		if isFileExist(filename) {
			return true
		}
	}

	return false
}
