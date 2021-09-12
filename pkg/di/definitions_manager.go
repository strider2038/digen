package di

import (
	"os"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

type DefinitionsManager struct {
	container *RootContainerDefinition
	workDir   string
}

func NewDefinitionsManager(container *RootContainerDefinition, workDir string) *DefinitionsManager {
	return &DefinitionsManager{
		container: container,
		workDir:   workDir,
	}
}

func (m *DefinitionsManager) Generate() ([]*File, error) {
	servicesByFiles := m.getServicesByFiles()

	files := make([]*File, 0, len(servicesByFiles))

	for filename, services := range servicesByFiles {
		if m.isDefinitionsFileExist(filename) {
			continue
		}

		file := NewFileBuilder(filename, "definitions", DefinitionsPackage)

		for _, service := range services {
			file.AddImport(m.container.GetImport(service))

			err := factoryTemplate.Execute(file, templateParameters{
				ServiceTitle: service.Title(),
				ServiceType:  service.Type.String(),
			})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate factory for %s", service.Name)
			}
		}

		files = append(files, file.GetFile())
	}

	return files, nil
}

func (m *DefinitionsManager) isDefinitionsFileExist(filename string) bool {
	_, err := os.Stat(m.workDir + "/" + packageDirs[DefinitionsPackage] + "/" + filename)

	return err == nil
}

func (m *DefinitionsManager) getServicesByFiles() map[string][]*ServiceDefinition {
	servicesByFiles := make(map[string][]*ServiceDefinition)

	for _, service := range m.container.Services {
		if service.IsExternal || service.IsRequired {
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

			servicesByFiles[filename] = append(servicesByFiles[filename], service)
		}
	}

	return servicesByFiles
}
