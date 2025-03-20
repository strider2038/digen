package di

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/spf13/afero"
)

type FactoriesGenerator struct {
	fs          afero.Fs
	fileLocator FileLocator
	container   *RootContainerDefinition
	params      GenerationParameters
}

func NewFactoriesGenerator(
	fs afero.Fs,
	fileLocator FileLocator,
	container *RootContainerDefinition,
	params GenerationParameters,
) *FactoriesGenerator {
	return &FactoriesGenerator{
		fs:          fs,
		fileLocator: fileLocator,
		container:   container,
		params:      params,
	}
}

func (g *FactoriesGenerator) Generate() ([]*File, error) {
	servicesByFiles := g.getServicesByFiles()

	files := make([]*File, 0, len(servicesByFiles))

	for filename, services := range servicesByFiles {
		var file *File
		var err error

		if isFileExist(g.fs, filename) {
			file, err = g.generateAppendFile(filename, services)
		} else {
			file, err = g.generateNewFile(filename, services)
		}
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

func (g *FactoriesGenerator) generateNewFile(filename string, services []*ServiceDefinition) (*File, error) {
	file := NewFileBuilder(filename, "factories")
	file.AddImportAliases(g.container.Imports)

	for _, service := range services {
		returnCode := make([]jen.Code, 0, 2)
		returnCode = append(returnCode, jen.Do(g.container.Type(service.Type)))
		if g.params.Factories.ReturnError() {
			returnCode = append(returnCode, jen.Error())
		}

		file.Add(
			jen.Line(),
			jen.Func().Id("Create"+strings.Title(service.Prefix)+service.Title()).
				Params(
					jen.Id("ctx").Qual("context", "Context"),
					jen.Id("c").Qual(g.params.packageName(LookupPackage), "Container"),
				).
				Params(returnCode...).
				Block(
					jen.Panic(jen.Lit("not implemented")),
				),
		)
	}

	return file.GetFile()
}

func (g *FactoriesGenerator) generateAppendFile(filename string, services []*ServiceDefinition) (*File, error) {
	var content bytes.Buffer

	for _, service := range services {
		factoryName := strings.Title(service.Prefix) + service.Title()
		if _, exists := g.container.Factories[factoryName]; exists {
			continue
		}

		returnCode := make([]jen.Code, 0, 2)
		returnCode = append(returnCode, jen.Do(g.container.Type(service.Type)))
		if g.params.Factories.ReturnError() {
			returnCode = append(returnCode, jen.Error())
		}

		content.WriteString("\n")
		content.WriteString(fmt.Sprintf("%#v",
			jen.Func().Id("Create"+factoryName).
				Params(
					jen.Id("ctx").Qual("context", "Context"),
					jen.Id("c").Qual(g.params.packageName(LookupPackage), "Container"),
				).
				Params(returnCode...).
				Block(jen.Panic(jen.Lit("not implemented"))),
		))
		content.WriteString("\n")
	}

	return &File{
		Name:    filename,
		Content: content.Bytes(),
		Append:  true,
	}, nil
}

func (g *FactoriesGenerator) getServicesByFiles() map[string][]*ServiceDefinition {
	servicesByFiles := make(map[string][]*ServiceDefinition)

	for _, service := range g.container.Services {
		if service.IsRequired {
			continue
		}

		filename := g.fileLocator.GetFactoryFilePath(service, "container.go")
		servicesByFiles[filename] = append(servicesByFiles[filename], service)
	}

	for _, container := range g.container.Containers {
		defaultFilename := strcase.ToSnake(container.Name) + ".go"

		for _, service := range container.Services {
			if service.IsRequired {
				continue
			}

			filename := g.fileLocator.GetFactoryFilePath(service, defaultFilename)
			servicesByFiles[filename] = append(servicesByFiles[filename], service)
		}
	}

	return servicesByFiles
}
