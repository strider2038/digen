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
	fs        afero.Fs
	container *RootContainerDefinition
	workDir   string
	params    GenerationParameters
}

func NewFactoriesGenerator(
	fs afero.Fs,
	container *RootContainerDefinition,
	workDir string,
	params GenerationParameters,
) *FactoriesGenerator {
	return &FactoriesGenerator{
		fs:        fs,
		container: container,
		workDir:   workDir,
		params:    params,
	}
}

func (g *FactoriesGenerator) Generate() ([]*File, error) {
	servicesByFiles := g.getServicesByFiles()

	files := make([]*File, 0, len(servicesByFiles))

	for filename, services := range servicesByFiles {
		var file *File
		var err error

		if g.isFactoryFileExist(filename) {
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
	file := NewFileBuilder(filename, "factories", FactoriesPackage)
	file.AddImportAliases(g.container.Imports)

	for _, service := range services {
		file.Add(
			jen.Line(),
			jen.Func().Id("Create"+strings.Title(service.Prefix)+service.Title()).
				Params(
					jen.Id("ctx").Qual("context", "Context"),
					jen.Id("c").Qual(g.params.packageName(LookupPackage), "Container"),
				).
				Params(
					jen.Do(g.container.Type(service.Type)),
					jen.Error(),
				).
				Block(jen.Panic(jen.Lit("not implemented"))),
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

		content.WriteString("\n")
		content.WriteString(fmt.Sprintf("%#v",
			jen.Func().Id("Create"+factoryName).
				Params(
					jen.Id("ctx").Qual("context", "Context"),
					jen.Id("c").Qual(g.params.packageName(LookupPackage), "Container"),
				).
				Params(
					jen.Do(g.container.Type(service.Type)),
					jen.Error(),
				).
				Block(jen.Panic(jen.Lit("not implemented"))),
		))
		content.WriteString("\n")
	}

	return &File{
		Package: FactoriesPackage,
		Name:    filename,
		Content: content.Bytes(),
		Append:  true,
	}, nil
}

func (g *FactoriesGenerator) isFactoryFileExist(filename string) bool {
	return isFileExist(g.fs, g.workDir+"/"+packageDirs[FactoriesPackage]+"/"+filename)
}

func (g *FactoriesGenerator) getServicesByFiles() map[string][]*ServiceDefinition {
	servicesByFiles := make(map[string][]*ServiceDefinition)

	for _, service := range g.container.Services {
		if service.IsRequired {
			continue
		}
		if service.FactoryFileName != "" {
			servicesByFiles[service.FactoryFileName] = append(servicesByFiles[service.FactoryFileName], service)
			continue
		}

		servicesByFiles["container.go"] = append(servicesByFiles["container.go"], service)
	}

	for _, container := range g.container.Containers {
		filename := strcase.ToSnake(container.Name) + ".go"

		for _, service := range container.Services {
			if service.IsRequired {
				continue
			}
			if service.FactoryFileName != "" {
				servicesByFiles[service.FactoryFileName] = append(servicesByFiles[service.FactoryFileName], service)
				continue
			}

			servicesByFiles[filename] = append(servicesByFiles[filename], service)
		}
	}

	return servicesByFiles
}
