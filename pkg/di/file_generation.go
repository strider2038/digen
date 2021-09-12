package di

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type GenerationParameters struct {
	RootPackage string
}

func (params GenerationParameters) rootPackageName() string {
	path := strings.Split(params.RootPackage, "/")
	if len(path) == 0 {
		return ""
	}
	return path[len(path)-1]
}

func (params *GenerationParameters) packageName(packageType PackageType) string {
	return strconv.Quote(params.RootPackage + "/" + packageDirs[packageType])
}

func GenerateFiles(container *RootContainerDefinition, params GenerationParameters) ([]*File, error) {
	files := make([]*File, 0)

	generators := [...]func(container *RootContainerDefinition, params GenerationParameters) (*File, error){
		func(container *RootContainerDefinition, params GenerationParameters) (*File, error) {
			generator := NewContainerGenerator(container, params)

			return generator.Generate()
		},
		generateDefinitionsContractsFile,
		func(container *RootContainerDefinition, params GenerationParameters) (*File, error) {
			generator := NewPublicContainerGenerator(container, params)

			return generator.Generate()
		},
	}

	for _, generate := range generators {
		file, err := generate(container, params)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to generate file")
		}

		files = append(files, file)
	}

	return files, nil
}

func GenerateContainerFile(params GenerationParameters) (*File, error) {
	var buffer bytes.Buffer

	err := internalContainerTemplate.Execute(&buffer, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate internal container")
	}

	return &File{
		Package: InternalPackage,
		Name:    "container.go",
		Content: buffer.Bytes(),
	}, nil
}

func generateDefinitionsContractsFile(container *RootContainerDefinition, params GenerationParameters) (*File, error) {
	file := NewFileBuilder("contracts.go", "definitions", DefinitionsPackage)

	file.WriteString("\ntype Container interface {\n")
	file.WriteString("\tSetError(err error)\n\n")
	for _, service := range container.Services {
		file.AddImport(container.GetImport(service))
		file.WriteString("\t" + service.Title() + "() " + service.Type.String() + "\n")
	}
	if len(container.Containers) > 0 {
		file.WriteString("\n")
		for _, attachedContainer := range container.Containers {
			file.WriteString("\t" + attachedContainer.Title() + "() " + attachedContainer.Type.Name + "\n")
		}
	}
	file.WriteString("}\n")

	for _, attachedContainer := range container.Containers {
		file.WriteString("\ntype " + attachedContainer.Type.Name + " interface {\n")
		for _, service := range attachedContainer.Services {
			file.AddImport(container.GetImport(service))
			file.WriteString("\t" + service.Title() + "() " + service.Type.String() + "\n")
		}
		file.WriteString("}\n")
	}

	return file.GetFile(), nil
}
