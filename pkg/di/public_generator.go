package di

import (
	"bytes"

	"github.com/pkg/errors"
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
			ContainerName: g.container.Name,
			ServiceName:   service.Name,
			ServiceTitle:  service.Title(),
			ServiceType:   service.Type.String(),
		}

		needsImport := false

		if service.IsPublic {
			needsImport = true
			err := publicGetterTemplate.Execute(&g.methods, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate getter for %s", service.Name)
			}
		}
		if service.HasSetter {
			needsImport = true
			err := publicSetterTemplate.Execute(&g.methods, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate setter for %s", service.Name)
			}
		}
		if service.IsRequired {
			needsImport = true
			err := argumentTemplate.Execute(&g.arguments, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate argument for %s", service.Name)
			}
			err = argumentSetterTemplate.Execute(&g.argumentSetters, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate argument setter for %s", service.Name)
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
				ServicePath:   attachedContainer.Title() + "().(" + attachedContainer.Type.String() + ").",
				ServiceName:   service.Name,
				ServiceTitle:  service.Title(),
				ServiceType:   service.Type.String(),
			}

			needsImport := false

			if service.IsPublic {
				needsImport = true
				err := publicGetterTemplate.Execute(&g.methods, parameters)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to generate getter for %s", service.Name)
				}
			}
			if service.HasSetter {
				needsImport = true
				err := publicSetterTemplate.Execute(&g.methods, parameters)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to generate setter for %s", service.Name)
				}
			}
			if service.IsRequired {
				needsImport = true
				err := argumentTemplate.Execute(&g.arguments, parameters)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to generate argument for %s", service.Name)
				}
				err = argumentSetterTemplate.Execute(&g.argumentSetters, parameters)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to generate argument setter for %s", service.Name)
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
		return nil, errors.Wrap(err, "failed to generate public container")
	}

	return g.file.GetFile(), nil
}
