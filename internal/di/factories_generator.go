package di

import (
	"os"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

type FactoriesGenerator struct {
	container *RootContainerDefinition
	workDir   string
	params    GenerationParameters
}

func NewFactoriesGenerator(
	container *RootContainerDefinition,
	workDir string,
	params GenerationParameters,
) *FactoriesGenerator {
	return &FactoriesGenerator{
		container: container,
		workDir:   workDir,
		params:    params,
	}
}

func (g *FactoriesGenerator) Generate() ([]*File, error) {
	servicesByFiles := g.getServicesByFiles()

	files := make([]*File, 0, len(servicesByFiles))

	for filename, services := range servicesByFiles {
		if g.isFactoryFileExist(filename) {
			continue
		}

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
					Do(g.container.Type(service.Type)).
					Block(jen.Panic(jen.Lit("not implemented"))),
			)
		}

		content, err := file.GetFile()
		if err != nil {
			return nil, err
		}
		files = append(files, content)
	}

	return files, nil
}

func (g *FactoriesGenerator) isFactoryFileExist(filename string) bool {
	_, err := os.Stat(g.workDir + "/" + packageDirs[FactoriesPackage] + "/" + filename)

	return err == nil
}

func (g *FactoriesGenerator) getServicesByFiles() map[string][]*ServiceDefinition {
	servicesByFiles := make(map[string][]*ServiceDefinition)

	for _, service := range g.container.Services {
		if service.IsExternal || service.IsRequired {
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
			if service.IsExternal || service.IsRequired {
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
