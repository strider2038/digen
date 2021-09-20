package di

import (
	"bytes"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"
)

type Generator struct {
	BaseDir        string
	ModulePath     string
	ConfigFilename string
	Logger         Logger

	Version   string
	BuildTime string
}

func (g Generator) RootPackage() string {
	return g.ModulePath + "/" + g.BaseDir
}

func (g *Generator) Generate() error {
	err := g.init()
	if err != nil {
		return err
	}

	container, err := ParseContainerFromFile(g.ConfigFilename)
	if err != nil {
		return err
	}

	g.Logger.Info("DI container", g.ConfigFilename, "successfully parsed")

	err = g.generateFiles(container)
	if err != nil {
		return err
	}

	err = g.generateDefinitionsFiles(container)
	if err != nil {
		return err
	}

	g.Logger.Success("Generation completed at dir", g.BaseDir)

	return nil
}

func (g *Generator) Initialize() error {
	err := g.init()
	if err != nil {
		return err
	}

	params := GenerationParameters{}
	file, err := GenerateContainerFile(params)
	if err != nil {
		return err
	}

	writer := NewWriter(g.BaseDir)

	err = writer.WriteFile(file)
	if err != nil {
		return err
	}

	g.Logger.Success("Init completed: file", file.Name, "generated")

	return nil
}

func (g *Generator) init() error {
	if g.BaseDir == "" {
		g.BaseDir = "."
	}
	if g.ConfigFilename == "" {
		g.ConfigFilename = g.BaseDir + "/internal/_config.go"
	}

	mod, err := os.ReadFile("go.mod")
	if err != nil {
		return errors.Wrap(err, "failed to read go.mod file")
	}
	path := modfile.ModulePath(mod)
	if path == "" {
		return errors.WithStack(errMissingModule)
	}

	g.ModulePath = path

	if g.Logger == nil {
		g.Logger = nilLogger{}
	}

	return nil
}

func (g *Generator) generateFiles(container *RootContainerDefinition) error {
	params := GenerationParameters{
		RootPackage: g.RootPackage(),
	}

	files, err := GenerateFiles(container, params)
	if err != nil {
		return err
	}

	writer := NewWriter(g.BaseDir)
	writer.Overwrite = true
	writer.Heading, err = g.generateHeading()
	if err != nil {
		return err
	}

	for _, file := range files {
		err = writer.WriteFile(file)
		if err != nil {
			return err
		}
		g.Logger.Info("File", file.Name, "generated")
	}

	return nil
}

func (g *Generator) generateDefinitionsFiles(container *RootContainerDefinition) error {
	manager := NewDefinitionsManager(container, g.BaseDir)
	files, err := manager.Generate()
	if err != nil {
		return err
	}

	for _, file := range files {
		writer := NewWriter(g.BaseDir)
		err = writer.WriteFile(file)
		if err != nil {
			return err
		}

		g.Logger.Info("Definitions file", file.Name, "generated")
	}

	return nil
}

func (g *Generator) isDefinitionsFileExist() bool {
	filename := g.BaseDir + "/" + packageDirs[DefinitionsPackage] + "/definitions.go"

	return isFileExist(filename)
}

func (g *Generator) generateHeading() ([]byte, error) {
	var heading bytes.Buffer
	err := headingTemplate.Execute(&heading, g)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate heading")
	}

	return heading.Bytes(), nil
}
