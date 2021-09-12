package di

import (
	"fmt"

	"github.com/pkg/errors"
)

type ContainerGenerator struct {
	container *RootContainerDefinition
	params    GenerationParameters

	file *FileBuilder
}

func NewContainerGenerator(container *RootContainerDefinition, params GenerationParameters) *ContainerGenerator {
	return &ContainerGenerator{
		container: container,
		params:    params,
		file:      NewFileBuilder("generated.go", container.Package, InternalPackage),
	}
}

func (g *ContainerGenerator) Generate() (*File, error) {
	generators := [...]func() error{
		g.generateConstructor,
		g.generateGetters,
		g.generateSetters,
		g.generateClosers,
	}

	for _, generate := range generators {
		err := generate()
		if err != nil {
			return nil, err
		}
	}

	return g.file.GetFile(), nil
}

func (g *ContainerGenerator) generateConstructor() error {
	g.writeLine("\nfunc NewContainer() *Container {")
	g.writeLine("\tc := &Container{}")

	for _, definition := range g.container.Containers {
		g.writeLine(fmt.Sprintf("\tc.%s = &%s{Container: c}", definition.Name, definition.Type.Name))
	}

	g.writeLine("\n\treturn c")
	g.writeLine("}")

	return nil
}

func (g *ContainerGenerator) generateGetters() error {
	err := g.writeServiceGetters(g.container.Services, g.container.Name)
	if err != nil {
		return err
	}

	for _, container := range g.container.Containers {
		g.importDefinitions()

		err = g.writeContainerGetter(container)
		if err != nil {
			return err
		}

		err = g.writeServiceGetters(container.Services, container.Type.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *ContainerGenerator) writeContainerGetter(container *ContainerDefinition) error {
	parameters := attachedContainerTemplateParameters{
		ContainerName:                   g.container.Name,
		AttachedContainerName:           container.Name,
		AttachedContainerTitle:          container.Title(),
		AttachedContainerDefinitionType: "definitions." + container.Type.Name,
	}
	err := attachedContainerGetterTemplate.Execute(g.file, parameters)
	if err != nil {
		return errors.Wrapf(err, "failed to generate getter for container %s", container.Name)
	}

	return nil
}

func (g *ContainerGenerator) writeServiceGetters(services []*ServiceDefinition, containerName string) error {
	for _, service := range services {
		g.importService(service)

		parameters := templateParameters{
			ContainerName: containerName,
			ServiceName:   service.Name,
			ServiceTitle:  service.Title(),
			ServiceType:   service.Type.String(),
			HasDefinition: !service.IsRequired,
			PanicOnNil:    service.IsExternal,
		}
		err := g.writeGetter(parameters, service)
		if err != nil {
			return err
		}

		if parameters.HasDefinition && !parameters.PanicOnNil {
			g.importDefinitions()
		}
	}

	return nil
}

func (g *ContainerGenerator) generateSetters() error {
	for _, service := range g.container.Services {
		err := g.generateSetter(g.container.Name, service)
		if err != nil {
			return err
		}
	}

	for _, attachedContainer := range g.container.Containers {
		for _, service := range attachedContainer.Services {
			err := g.generateSetter(attachedContainer.Type.Name, service)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *ContainerGenerator) generateClosers() error {
	g.write("\nfunc (c *Container) Close() {")
	for _, service := range g.container.Services {
		if service.HasCloser {
			err := closerTemplate.Execute(g.file, templateParameters{ServiceName: service.Name})
			if err != nil {
				return errors.Wrapf(err, "failed to generate closer for %s", service.Name)
			}
		}
	}

	for _, attachedContainer := range g.container.Containers {
		for _, service := range attachedContainer.Services {
			if service.HasCloser {
				err := closerTemplate.Execute(g.file, templateParameters{
					ServiceName: attachedContainer.Name + "." + service.Name,
				})
				if err != nil {
					return errors.Wrapf(err, "failed to generate closer for %s", service.Name)
				}
			}
		}
	}
	g.writeLine("}")

	return nil
}

func (g *ContainerGenerator) generateSetter(containerName string, service *ServiceDefinition) error {
	if service.HasSetter || service.IsExternal || service.IsRequired {
		g.importService(service)
		err := setterTemplate.Execute(g.file, templateParameters{
			ContainerName: containerName,
			ServiceName:   service.Name,
			ServiceTitle:  service.Title(),
			ServiceType:   service.Type.String(),
		})
		if err != nil {
			return errors.Wrapf(err, "failed to generate setter for %s", service.Name)
		}
	}

	return nil
}

func (g *ContainerGenerator) write(s string) {
	g.file.WriteString(s)
}

func (g *ContainerGenerator) writeLine(s string) {
	g.write(s)
	g.newLine()
}

func (g *ContainerGenerator) newLine() {
	g.write("\n")
}

func (g *ContainerGenerator) writeGetter(parameters templateParameters, service *ServiceDefinition) error {
	err := getterTemplate.Execute(g.file, parameters)
	if err != nil {
		return errors.Wrapf(err, "failed to generate getter for %s", service.Name)
	}
	return nil
}

func (g *ContainerGenerator) importService(service *ServiceDefinition) {
	g.file.AddImport(g.container.GetImport(service))
}

func (g *ContainerGenerator) importDefinitions() {
	g.file.AddImport(g.params.packageName(DefinitionsPackage))
}
