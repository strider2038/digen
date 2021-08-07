package console

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/strider2038/digen"
)

func Generate(options *Options) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	container, err := digen.ParseContainerFromFile(config.ContainerFilename)
	if err != nil {
		return err
	}

	pterm.Info.Println("DI container", config.ContainerFilename, "successfully parsed")

	err = generateFiles(container, config, options)
	if err != nil {
		return err
	}

	if !digen.IsDefinitionsFileExist(config.WorkDir) {
		err = generateDefinitionsFile(container, config)
		if err != nil {
			return err
		}
	}

	pterm.Success.Println("Generation completed at dir", config.WorkDir)

	return nil
}

func generateFiles(container *digen.ContainerDefinition, config *Configuration, options *Options) error {
	params := digen.GenerationParameters{
		RootPackage: config.RootPackage(),
	}

	files, err := digen.Generate(container, params)
	if err != nil {
		return err
	}

	writer := digen.NewWriter(config.WorkDir)
	writer.Overwrite = true
	writer.Heading, err = generateHeading(options)
	if err != nil {
		return err
	}

	for _, file := range files {
		err = writer.WriteFile(file)
		if err != nil {
			return err
		}
		pterm.Info.Println("File", file.Name, "generated")
	}

	return nil
}

func generateDefinitionsFile(container *digen.ContainerDefinition, config *Configuration) error {
	params := digen.GenerationParameters{
		RootPackage: config.RootPackage(),
	}
	file, err := digen.GenerateFactory(container, params)
	if err != nil {
		return err
	}

	writer := digen.NewWriter(config.WorkDir)
	err = writer.WriteFile(file)
	if err != nil {
		return err
	}

	pterm.Info.Println("File", file.Name, "generated")
	return nil
}

func generateHeading(options *Options) ([]byte, error) {
	var heading bytes.Buffer
	err := headingTemplate.Execute(&heading, options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate heading")
	}

	return heading.Bytes(), nil
}

var headingTemplate = template.Must(template.New("heading").Parse(`// Code generated by DIGEN; DO NOT EDIT.
// This file was generated by Dependency Injection Container Generator {{.Version}} (built at {{.BuildTime}}).
// See docs at https://github.com/strider2038/digen.

`))
