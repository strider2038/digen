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
	BaseDir    string
	ModulePath string
	Params     GenerationParameters

	FS          afero.Fs
	Logger      Logger
	FileLocator FileLocator
}

func (g *Generator) RootPackage() string {
	return g.ModulePath + "/" + g.BaseDir
}

func (g *Generator) Initialize() error {
	if err := g.init(); err != nil {
		return err
	}

	file := &File{
		Name:    g.FileLocator.GetPackageFilePath(DefinitionsPackage, "container.go"),
		Content: []byte(definitionsContainerFileSkeleton),
	}

	writer := NewWriter(g.FS)
	if err := writer.WriteFile(file); err != nil {
		if errors.Is(err, ErrFileAlreadyExists) {
			g.Logger.Warning("init skipped: file", file.Name, "already exists")

			return nil
		}

		return err
	}

	g.Logger.Success("init completed: file", file.Name, "generated")

	return nil
}

func (g *Generator) Generate() error {
	if err := g.init(); err != nil {
		return err
	}

	container, err := g.parseDefinitionsFromFile(g.BaseDir + "/" + definitionsFile)
	if err != nil {
		return errors.Errorf("parse definitions file: %w", err)
	}
	g.Logger.Info("service definitions parsed from:", definitionsFile)

	factories, err := g.parseFactories(container)
	if err != nil {
		return errors.Errorf("parse factories: %w", err)
	}
	if len(factories.Factories) > 0 {
		container.Factories = factories.Factories
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

	g.FileLocator = FileLocator{
		ContainerDir: g.BaseDir,
		ModulePath:   g.ModulePath,
	}

	g.Params = g.Params.Defaults()
	g.Params.RootPackage = g.RootPackage()

	return nil
}

func (g *Generator) parseDefinitionsFromFile(filename string) (*RootContainerDefinition, error) {
	parser := NewDefinitionsParser(g.FS, g.Logger)

	return parser.ParseFile(filename)
}

func (g *Generator) parseFactories(container *RootContainerDefinition) (*FactoryDefinitions, error) {
	dirs := []string{g.BaseDir + "/" + factoriesDir}

	dirVisited := make(map[string]struct{})

	for _, service := range container.Services {
		if service.FactoryPackage != "" {
			dir := g.FileLocator.GetPathByPackage(service.FactoryPackage)
			if _, visited := dirVisited[dir]; !visited {
				dirs = append(dirs, dir)
				dirVisited[dir] = struct{}{}
			}
		}
	}
	for _, c := range container.Containers {
		for _, service := range c.Services {
			if service.FactoryPackage != "" {
				dir := g.FileLocator.GetPathByPackage(service.FactoryPackage)
				if _, visited := dirVisited[dir]; !visited {
					dirs = append(dirs, dir)
					dirVisited[dir] = struct{}{}
				}
			}
		}
	}

	return parseFactoriesFromDirs(g.FS, g.Logger, dirs...)
}

func (g *Generator) generateContainerFiles(container *RootContainerDefinition) error {
	files, err := NewFileGenerator(g.FileLocator, container, g.Params).GenerateFiles()
	if err != nil {
		return err
	}

	writer := NewWriter(g.FS)
	writer.Overwrite = true

	for _, file := range files {
		err = writer.WriteFile(file)
		if err != nil {
			return err
		}
		g.Logger.Info("file", file.Name, "generated")
	}

	return nil
}

func (g *Generator) generateFactoriesFiles(container *RootContainerDefinition) error {
	generator := NewFactoriesGenerator(g.FS, g.FileLocator, container, g.Params)
	files, err := generator.Generate()
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsEmpty() {
			continue
		}
		writer := NewWriter(g.FS)
		writer.Append = file.Append
		err = writer.WriteFile(file)
		if err != nil {
			return err
		}

		action := "generated"
		if writer.Append {
			action = "updated"
		}
		g.Logger.Info("factories file", file.Name, action)
	}

	return nil
}

func (g *Generator) generateUtils() error {
	heading := []byte(fmt.Sprintf(headingTemplate, g.Params.Version))

	file := &File{
		Name:    g.FileLocator.GetPackageFilePath(InternalPackage, "bitset.go"),
		Content: slices.Concat(heading, []byte(bitsetSkeleton)),
	}

	writer := NewWriter(g.FS)
	writer.Overwrite = true
	if err := writer.WriteFile(file); err != nil {
		return err
	}

	g.Logger.Info("file", file.Name, "generated")

	return nil
}

func (g *Generator) generateReadmeFile() error {
	file := &File{
		Name:    g.FileLocator.GetContainerFilePath("README.md"),
		Content: []byte(readmeTemplate),
	}

	writer := NewWriter(g.FS)
	writer.Overwrite = true
	if err := writer.WriteFile(file); err != nil {
		return err
	}

	g.Logger.Info("readme file", file.Name, "generated")

	return nil
}
