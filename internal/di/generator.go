package di

import (
	"bytes"
	"os"

	"github.com/muonsoft/errors"
	"golang.org/x/mod/modfile"
)

const (
	definitionsFile = "internal/definitions/container.go"
)

type Generator struct {
	BaseDir    string
	ModulePath string
	Logger     Logger

	Version   string
	BuildTime string
}

func (g *Generator) RootPackage() string {
	return g.ModulePath + "/" + g.BaseDir
}

func (g *Generator) Initialize() error {
	if err := g.init(); err != nil {
		return err
	}

	file, err := GenerateDefinitionsContainerFile()
	if err != nil {
		return err
	}

	writer := NewWriter(g.BaseDir)
	if err := writer.WriteFile(file); err != nil {
		if errors.Is(err, ErrFileAlreadyExists) {
			g.Logger.Warning("init skipped: file", file.Path(), "already exists")

			return nil
		}

		return err
	}

	g.Logger.Success("init completed: file", file.Path(), "generated")

	return nil
}

func (g *Generator) Generate() error {
	if err := g.init(); err != nil {
		return err
	}

	container, err := ParseDefinitionsFromFile(g.BaseDir + "/" + definitionsFile)
	if err != nil {
		return err
	}

	g.Logger.Info("definitions container", definitionsFile, "successfully parsed")

	if err := g.generateContainerFiles(container); err != nil {
		return err
	}
	if err := g.generateFactoriesFiles(container); err != nil {
		return err
	}
	if err := g.generateReadmeFile(); err != nil {
		return err
	}

	g.Logger.Success("generation completed at dir", g.BaseDir)

	return nil
}

func (g *Generator) init() error {
	if g.BaseDir == "" {
		g.BaseDir = "."
	}

	mod, err := os.ReadFile("go.mod")
	if err != nil {
		return errors.Errorf("read go.mod file: %w", err)
	}
	path := modfile.ModulePath(mod)
	if path == "" {
		return errors.Wrap(errMissingModule)
	}

	g.ModulePath = path

	if g.Logger == nil {
		g.Logger = nilLogger{}
	}

	return nil
}

func (g *Generator) generateContainerFiles(container *RootContainerDefinition) error {
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
		g.Logger.Info("file", file.Path(), "generated")
	}

	return nil
}

func (g *Generator) generateFactoriesFiles(container *RootContainerDefinition) error {
	params := GenerationParameters{
		RootPackage: g.RootPackage(),
	}
	manager := NewFactoriesGenerator(container, g.BaseDir, params)
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

		g.Logger.Info("factories file", file.Path(), "generated")
	}

	return nil
}

func (g *Generator) generateReadmeFile() error {
	file := &File{
		Name:    "README.md",
		Content: []byte(readmeTemplate),
	}

	writer := NewWriter(g.BaseDir)
	writer.Overwrite = true
	if err := writer.WriteFile(file); err != nil {
		return err
	}

	g.Logger.Info("readme file", file.Path(), "generated")

	return nil
}

func (g *Generator) generateHeading() ([]byte, error) {
	var heading bytes.Buffer
	err := headingTemplate.Execute(&heading, g)
	if err != nil {
		return nil, errors.Errorf("generate heading: %w", err)
	}

	return heading.Bytes(), nil
}
