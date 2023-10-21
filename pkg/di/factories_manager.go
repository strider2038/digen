package di

import (
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/muonsoft/errors"
)

type FactoriesManager struct {
	container *RootContainerDefinition
	workDir   string
	params    GenerationParameters
}

func NewFactoriesManager(
	container *RootContainerDefinition,
	workDir string,
	params GenerationParameters,
) *FactoriesManager {
	return &FactoriesManager{
		container: container,
		workDir:   workDir,
		params:    params,
	}
}

func (m *FactoriesManager) Generate() ([]*File, error) {
	servicesByFiles := m.getServicesByFiles()

	files := make([]*File, 0, len(servicesByFiles))

	for filename, services := range servicesByFiles {
		if m.isFactoryFileExist(filename) {
			continue
		}

		file := NewFileBuilder(filename, "factories", FactoriesPackage)
		file.AddImport(`"context"`)
		file.AddImport(m.params.packageName(LookupPackage))

		for _, service := range services {
			file.AddImport(m.container.GetImport(service))

			err := factoryFuncTemplate.Execute(file, templateParameters{
				ServicePrefix: strings.Title(service.Prefix),
				ServiceTitle:  service.Title(),
				ServiceType:   service.Type.String(),
			})
			if err != nil {
				return nil, errors.Errorf("generate factory for %s: %w", service.Name, err)
			}
		}

		files = append(files, file.GetFile())
	}

	return files, nil
}

func (m *FactoriesManager) isFactoryFileExist(filename string) bool {
	_, err := os.Stat(m.workDir + "/" + packageDirs[FactoriesPackage] + "/" + filename)

	return err == nil
}

func (m *FactoriesManager) getServicesByFiles() map[string][]*ServiceDefinition {
	servicesByFiles := make(map[string][]*ServiceDefinition)

	for _, service := range m.container.Services {
		if service.IsExternal || service.IsRequired {
			continue
		}
		if service.FactoryFileName != "" {
			servicesByFiles[service.FactoryFileName] = append(servicesByFiles[service.FactoryFileName], service)
			continue
		}

		servicesByFiles["container.go"] = append(servicesByFiles["container.go"], service)
	}

	for _, container := range m.container.Containers {
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
