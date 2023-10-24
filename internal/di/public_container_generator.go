package di

import (
	"bytes"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/muonsoft/errors"
)

type PublicContainerGenerator struct {
	container *RootContainerDefinition
	params    GenerationParameters

	file *FileBuilder

	methods         bytes.Buffer
	arguments       bytes.Buffer
	argumentSetters bytes.Buffer
}

func NewPublicContainerGenerator(
	container *RootContainerDefinition,
	params GenerationParameters,
) *PublicContainerGenerator {
	file := NewFileBuilder("container.go", params.rootPackageName(), PublicPackage)
	file.AddImport(`"sync"`)
	file.AddImport(params.packageName(InternalPackage))

	return &PublicContainerGenerator{
		container: container,
		params:    params,
		file:      file,
	}
}

func (g *PublicContainerGenerator) Generate() (*File, error) {
	for _, service := range g.container.Services {
		parameters := templateParameters{
			ContainerName:     g.container.Name,
			ServiceName:       strcase.ToLowerCamel(service.Name),
			ServiceTitle:      service.Title(),
			ServiceType:       service.Type.String(),
			ServicePublicName: strings.Title(service.PublicName),
		}

		needsImport := false

		if service.IsPublic {
			g.file.AddImport(`"context"`)
			needsImport = true
			err := publicGetterTemplate.Execute(&g.methods, parameters)
			if err != nil {
				return nil, errors.Errorf("generate getter for %s: %w", service.Name, err)
			}
		}
		if service.HasSetter {
			needsImport = true
			err := publicSetterTemplate.Execute(&g.methods, parameters)
			if err != nil {
				return nil, errors.Errorf("generate setter for %s: %w", service.Name, err)
			}
		}
		if service.IsRequired {
			needsImport = true
			err := argumentTemplate.Execute(&g.arguments, parameters)
			if err != nil {
				return nil, errors.Errorf("generate argument for %s: %w", service.Name, err)
			}
			err = argumentSetterTemplate.Execute(&g.argumentSetters, parameters)
			if err != nil {
				return nil, errors.Errorf("generate argument setter for %s: %w", service.Name, err)
			}
		}

		if needsImport {
			g.file.AddImport(g.container.GetImport(service))
		}
	}

	for _, attachedContainer := range g.container.Containers {
		for _, service := range attachedContainer.Services {
			parameters := templateParameters{
				ContainerName: g.container.Name,
				ServicePath:   attachedContainer.Title() + "().(*" + attachedContainer.Type.String() + ").",
				ServiceName:   service.Name,
				ServiceTitle:  service.Title(),
				ServiceType:   service.Type.String(),
			}

			needsImport := false

			if service.IsPublic {
				g.file.AddImport(`"context"`)
				needsImport = true
				err := publicGetterTemplate.Execute(&g.methods, parameters)
				if err != nil {
					return nil, errors.Errorf("generate getter for %s: %w", service.Name, err)
				}
			}
			if service.HasSetter {
				needsImport = true
				err := publicSetterTemplate.Execute(&g.methods, parameters)
				if err != nil {
					return nil, errors.Errorf("generate setter for %s: %w", service.Name, err)
				}
			}
			if service.IsRequired {
				needsImport = true
				err := argumentTemplate.Execute(&g.arguments, parameters)
				if err != nil {
					return nil, errors.Errorf("generate argument for %s: %w", service.Name, err)
				}
				err = argumentSetterTemplate.Execute(&g.argumentSetters, parameters)
				if err != nil {
					return nil, errors.Errorf("generate argument setter for %s: %w", service.Name, err)
				}
			}

			if needsImport {
				g.file.AddImport(g.container.GetImport(service))
			}
		}
	}

	if g.arguments.Len() == 0 {
		g.arguments.WriteString("injectors ...Injector")
	} else {
		g.arguments.WriteString("\n\tinjectors ...Injector,\n")
		g.argumentSetters.WriteString("\n")
	}

	err := publicContainerTemplate.Execute(g.file, containerTemplateParameters{
		ContainerArguments:       g.arguments.String(),
		ContainerArgumentSetters: g.argumentSetters.String(),
		ContainerMethods:         g.methods.String(),
	})
	if err != nil {
		return nil, errors.Errorf("generate public container: %w", err)
	}

	return g.file.GetFile(), nil
}
