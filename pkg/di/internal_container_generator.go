package di

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/muonsoft/errors"
)

type InternalContainerGenerator struct {
	container *RootContainerDefinition
	params    GenerationParameters

	file *FileBuilder
}

func NewInternalContainerGenerator(container *RootContainerDefinition, params GenerationParameters) *InternalContainerGenerator {
	return &InternalContainerGenerator{
		container: container,
		params:    params,
		file:      NewFileBuilder("container.go", "internal", InternalPackage),
	}
}

func (g *InternalContainerGenerator) Generate() (*File, error) {
	generators := [...]func() error{
		g.generateRootContainer,
		g.generateContainers,
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

func (g *InternalContainerGenerator) generateRootContainer() error {
	g.file.AddImport(`"context"`)

	var body bytes.Buffer
	for _, service := range g.container.Services {
		g.importService(service)
		body.WriteString(fmt.Sprintf("\n\t%s %s", strcase.ToLowerCamel(service.Name), service.Type.String()))
	}

	if len(g.container.Containers) > 0 {
		g.file.AddImport(g.params.packageName(LookupPackage))
		body.WriteString("\n")
	}

	var constructor bytes.Buffer
	for _, container := range g.container.Containers {
		body.WriteString(fmt.Sprintf(
			"\n\t%s *%s",
			strcase.ToLowerCamel(container.Name), container.Type.Name,
		))
		constructor.WriteString(fmt.Sprintf(
			"\tc.%s = &%s{Container: c}\n",
			strcase.ToLowerCamel(container.Name), container.Type.Name),
		)
	}

	err := internalContainerTemplate.Execute(g.file, internalContainerTemplateParameters{
		ContainerBody:        body.String(),
		ContainerConstructor: constructor.String(),
	})
	if err != nil {
		return errors.Errorf("generate root container: %w", err)
	}

	return nil
}

func (g *InternalContainerGenerator) generateContainers() error {
	for _, container := range g.container.Containers {
		g.writeLine(fmt.Sprintf("\ntype %s struct {", container.Type.Name))
		g.writeLine("\t*Container\n")

		for _, service := range container.Services {
			g.importService(service)
			g.writeLine(fmt.Sprintf("\t%s %s", strcase.ToLowerCamel(service.Name), service.Type.String()))
		}

		g.writeLine("}")
	}

	return nil
}

func (g *InternalContainerGenerator) generateGetters() error {
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

func (g *InternalContainerGenerator) writeContainerGetter(container *ContainerDefinition) error {
	parameters := attachedContainerTemplateParameters{
		ContainerName:                   g.container.Name,
		AttachedContainerName:           strcase.ToLowerCamel(container.Name),
		AttachedContainerTitle:          container.Title(),
		AttachedContainerDefinitionType: "lookup." + container.Type.Name,
	}
	err := separateContainerGetterTemplate.Execute(g.file, parameters)
	if err != nil {
		return errors.Errorf("generate getter for container %s: %w", container.Name, err)
	}

	return nil
}

func (g *InternalContainerGenerator) writeServiceGetters(services []*ServiceDefinition, containerName string) error {
	for _, service := range services {
		g.importService(service)

		parameters := templateParameters{
			ContainerName:         containerName,
			ServicePrefix:         strings.Title(service.Prefix),
			ServiceName:           strcase.ToLowerCamel(service.Name),
			ServiceTitle:          service.Title(),
			ServiceType:           service.Type.String(),
			ServiceZeroComparison: service.Type.ZeroComparison(),
			HasDefinition:         !service.IsRequired,
			PanicOnNil:            service.IsExternal,
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

func (g *InternalContainerGenerator) generateSetters() error {
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

func (g *InternalContainerGenerator) generateClosers() error {
	g.write("\nfunc (c *Container) Close() {")
	for _, service := range g.container.Services {
		if service.HasCloser {
			err := closerTemplate.Execute(g.file, templateParameters{ServiceName: strcase.ToLowerCamel(service.Name)})
			if err != nil {
				return errors.Errorf("generate closer for %s: %w", service.Name, err)
			}
		}
	}

	for _, attachedContainer := range g.container.Containers {
		for _, service := range attachedContainer.Services {
			if service.HasCloser {
				err := closerTemplate.Execute(g.file, templateParameters{
					ServiceName: attachedContainer.Name + "." + strcase.ToLowerCamel(service.Name),
				})
				if err != nil {
					return errors.Errorf("generate closer for %s: %w", service.Name, err)
				}
			}
		}
	}
	g.writeLine("}")

	return nil
}

func (g *InternalContainerGenerator) generateSetter(containerName string, service *ServiceDefinition) error {
	if service.HasSetter || service.IsExternal || service.IsRequired {
		g.importService(service)
		err := setterTemplate.Execute(g.file, templateParameters{
			ContainerName: containerName,
			ServiceName:   strcase.ToLowerCamel(service.Name),
			ServiceTitle:  service.Title(),
			ServiceType:   service.Type.String(),
		})
		if err != nil {
			return errors.Errorf("generate setter for %s: %w", service.Name, err)
		}
	}

	return nil
}

func (g *InternalContainerGenerator) write(s string) {
	g.file.WriteString(s)
}

func (g *InternalContainerGenerator) writeLine(s string) {
	g.write(s)
	g.newLine()
}

func (g *InternalContainerGenerator) newLine() {
	g.write("\n")
}

func (g *InternalContainerGenerator) writeGetter(parameters templateParameters, service *ServiceDefinition) error {
	err := getterTemplate.Execute(g.file, parameters)
	if err != nil {
		return errors.Errorf("generate getter for %s: %w", service.Name, err)
	}
	return nil
}

func (g *InternalContainerGenerator) importService(service *ServiceDefinition) {
	g.file.AddImport(g.container.GetImport(service))
}

func (g *InternalContainerGenerator) importDefinitions() {
	g.file.AddImport(g.params.packageName(FactoriesPackage))
}
