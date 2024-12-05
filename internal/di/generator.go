package di

import (
	"bytes"
	"os"
	"slices"

	"github.com/muonsoft/errors"
	"golang.org/x/mod/modfile"
)

const (
	definitionsFile = "internal/definitions/container.go"
	factoriesDir    = "internal/factories"
)

type Generator struct {
	BaseDir       string
	ModulePath    string
	Logger        Logger
	ErrorHandling ErrorHandling

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

	file := GenerateDefinitionsContainerFile()

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
		return errors.Errorf("parse definitions file: %w", err)
	}
	g.Logger.Info("definitions container", definitionsFile, "successfully parsed")

	factories, err := ParseFactoriesFromDir(g.BaseDir + "/" + factoriesDir)
	if err != nil {
		return errors.Errorf("parse factories: %w", err)
	}
	if len(factories.Factories) > 0 {
		container.Factories = factories.Factories
		g.Logger.Info("factories", factoriesDir, "successfully parsed")
	}

	if err := g.generateContainerFiles(container); err != nil {
		return err
	}
	if err := g.generateFactoriesFiles(container); err != nil {
		return err
	}
	if err := g.generateUtils(); err != nil {
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
		RootPackage:   g.RootPackage(),
		ErrorHandling: g.ErrorHandling.Defaults(),
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
		RootPackage:   g.RootPackage(),
		ErrorHandling: g.ErrorHandling.Defaults(),
	}
	generator := NewFactoriesGenerator(container, g.BaseDir, params)
	files, err := generator.Generate()
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsEmpty() {
			continue
		}
		writer := NewWriter(g.BaseDir)
		writer.Append = file.Append
		err = writer.WriteFile(file)
		if err != nil {
			return err
		}

		action := "generated"
		if writer.Append {
			action = "updated"
		}
		g.Logger.Info("factories file", file.Path(), action)
	}

	return nil
}

func (g *Generator) generateUtils() error {
	heading, err := g.generateHeading()
	if err != nil {
		return err
	}

	file := &File{
		Package: InternalPackage,
		Name:    "bitset.go",
		Content: slices.Concat(heading, []byte(bitsetSkeleton)),
	}

	writer := NewWriter(g.BaseDir)
	writer.Overwrite = true
	if err := writer.WriteFile(file); err != nil {
		return err
	}

	g.Logger.Info("file", file.Path(), "generated")

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
