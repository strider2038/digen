package di

import (
	"fmt"
	"slices"

	"github.com/muonsoft/errors"
	"github.com/spf13/afero"
	"golang.org/x/mod/modfile"
)

const (
	definitionsFile = "internal/definitions/container.go"
	factoriesDir    = "internal/factories"
)

type Generator struct {
	BaseDir       string
	ModulePath    string
	ErrorHandling ErrorHandling

	FS     afero.Fs
	Logger Logger

	Version string

	params GenerationParameters
}

func (g *Generator) RootPackage() string {
	return g.ModulePath + "/" + g.BaseDir
}

func (g *Generator) Initialize() error {
	if err := g.init(); err != nil {
		return err
	}

	file := generateDefinitionsContainerFile()

	writer := NewWriter(g.FS, g.BaseDir)
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

	container, err := ParseDefinitionsFromFile(g.FS, g.BaseDir+"/"+definitionsFile)
	if err != nil {
		return errors.Errorf("parse definitions file: %w", err)
	}
	g.Logger.Info("definitions container", definitionsFile, "successfully parsed")

	factories, err := parseFactoriesFromDir(g.FS, g.BaseDir+"/"+factoriesDir)
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
	if g.FS == nil {
		g.FS = afero.NewOsFs()
	}

	if g.ModulePath == "" {
		mod, err := afero.ReadFile(g.FS, "go.mod")
		if err != nil {
			return errors.Errorf("read go.mod file: %w", err)
		}
		path := modfile.ModulePath(mod)
		if path == "" {
			return errors.Wrap(errMissingModule)
		}

		g.ModulePath = path
	}

	if g.Logger == nil {
		g.Logger = nilLogger{}
	}

	g.params.RootPackage = g.RootPackage()
	g.params.ErrorHandling = g.ErrorHandling.Defaults()
	g.params.Version = g.Version

	return nil
}

func (g *Generator) generateContainerFiles(container *RootContainerDefinition) error {
	files, err := NewFileGenerator(container, g.params).GenerateFiles()
	if err != nil {
		return err
	}

	writer := NewWriter(g.FS, g.BaseDir)
	writer.Overwrite = true

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
	generator := NewFactoriesGenerator(g.FS, container, g.BaseDir, g.params)
	files, err := generator.Generate()
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsEmpty() {
			continue
		}
		writer := NewWriter(g.FS, g.BaseDir)
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
	heading := []byte(fmt.Sprintf(headingTemplate, g.Version))

	file := &File{
		Package: InternalPackage,
		Name:    "bitset.go",
		Content: slices.Concat(heading, []byte(bitsetSkeleton)),
	}

	writer := NewWriter(g.FS, g.BaseDir)
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

	writer := NewWriter(g.FS, g.BaseDir)
	writer.Overwrite = true
	if err := writer.WriteFile(file); err != nil {
		return err
	}

	g.Logger.Info("readme file", file.Path(), "generated")

	return nil
}
