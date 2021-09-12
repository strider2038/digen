package digen

import (
	"bytes"
	"fmt"
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
		generateContainerConstructorFile,
		generateContainerGettersFile,
		generateContainerSettersFile,
		generateContainerCloseFile,
		generateDefinitionsContractsFile,
		generatePublicContainerFile,
	}

	for _, generate := range generators {
		file, err := generate(container, params)
		if errors.Is(err, errFileIgnored) {
			continue
		}
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

func GenerateFactory(container *RootContainerDefinition, params GenerationParameters) (*File, error) {
	file := NewFileBuilder("definitions.go", "definitions", DefinitionsPackage)

	for _, service := range container.Services {
		if service.IsExternal || service.IsRequired {
			continue
		}
		file.AddImport(container.GetImport(service))

		err := factoryTemplate.Execute(file, templateParameters{
			ServiceTitle: service.Title(),
			ServiceType:  service.Type.String(),
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate factory for %s", service.Name)
		}
	}

	return file.GetFile(), nil
}

func generateContainerConstructorFile(container *RootContainerDefinition, params GenerationParameters) (*File, error) {
	file := NewFileBuilder("container_new.go", container.Package, InternalPackage)

	file.WriteString("func NewContainer() *Container {\n")
	file.WriteString("\tc := &Container{}\n")

	for _, definition := range container.Containers {
		s := fmt.Sprintf("\tc.%s = &%s{Container: c}\n", definition.Name, definition.Type.Name)
		file.WriteString(s)
	}

	file.WriteString("\n\treturn c\n")
	file.WriteString("}\n")

	return file.GetFile(), nil
}

func generateContainerGettersFile(container *RootContainerDefinition, params GenerationParameters) (*File, error) {
	file := NewFileBuilder("container_get.go", container.Package, InternalPackage)
	needsDefinitions := false

	for _, service := range container.Services {
		file.AddImport(container.GetImport(service))

		parameters := templateParameters{
			ContainerName: container.Name,
			ServiceName:   service.Name,
			ServiceTitle:  service.Title(),
			ServiceType:   service.Type.String(),
			HasDefinition: !service.IsRequired,
			PanicOnNil:    service.IsExternal,
		}
		err := getterTemplate.Execute(file, parameters)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate getter for %s", service.Name)
		}

		if parameters.HasDefinition && !parameters.PanicOnNil {
			needsDefinitions = true
		}
	}

	for _, attachedContainer := range container.Containers {
		parameters := attachedContainerTemplateParameters{
			ContainerName:                   container.Name,
			AttachedContainerName:           attachedContainer.Name,
			AttachedContainerTitle:          attachedContainer.Title(),
			AttachedContainerDefinitionType: "definitions." + attachedContainer.Type.Name,
		}
		err := attachedContainerGetterTemplate.Execute(file, parameters)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate getter for container %s", attachedContainer.Name)
		}

		needsDefinitions = true

		for _, service := range attachedContainer.Services {
			file.AddImport(container.GetImport(service))

			parameters := templateParameters{
				ContainerName: attachedContainer.Type.Name,
				ServiceName:   service.Name,
				ServiceTitle:  service.Title(),
				ServiceType:   service.Type.String(),
				HasDefinition: !service.IsRequired,
				PanicOnNil:    service.IsExternal,
			}
			err := getterTemplate.Execute(file, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate getter for %s", service.Name)
			}
		}
	}

	if needsDefinitions {
		file.AddImport(params.packageName(DefinitionsPackage))
	}

	return file.GetFile(), nil
}

func generateContainerSettersFile(container *RootContainerDefinition, params GenerationParameters) (*File, error) {
	file := NewFileBuilder("container_set.go", container.Package, InternalPackage)

	for _, service := range container.Services {
		err := generateSetter(file, container, container.Name, service)
		if err != nil {
			return nil, err
		}
	}

	for _, attachedContainer := range container.Containers {
		for _, service := range attachedContainer.Services {
			err := generateSetter(file, container, attachedContainer.Type.Name, service)
			if err != nil {
				return nil, err
			}
		}
	}

	if file.IsEmpty() {
		return nil, errFileIgnored
	}

	return file.GetFile(), nil
}

func generateSetter(file *FileBuilder, container *RootContainerDefinition, containerName string, service *ServiceDefinition) error {
	if service.HasSetter || service.IsExternal || service.IsRequired {
		file.AddImport(container.GetImport(service))
		err := setterTemplate.Execute(file, templateParameters{
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

func generateContainerCloseFile(container *RootContainerDefinition, params GenerationParameters) (*File, error) {
	file := NewFileBuilder("container_close.go", container.Package, InternalPackage)

	file.WriteString("func (c *Container) Close() {")
	for _, service := range container.Services {
		if service.HasCloser {
			err := closerTemplate.Execute(file, templateParameters{ServiceName: service.Name})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate closer for %s", service.Name)
			}
		}

		for _, attachedContainer := range container.Containers {
			for _, service := range attachedContainer.Services {
				if service.HasCloser {
					err := closerTemplate.Execute(file, templateParameters{
						ServiceName: attachedContainer.Name + "." + service.Name,
					})
					if err != nil {
						return nil, errors.Wrapf(err, "failed to generate closer for %s", service.Name)
					}
				}
			}
		}
	}
	file.WriteString("}\n")

	return file.GetFile(), nil
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

func generatePublicContainerFile(container *RootContainerDefinition, params GenerationParameters) (*File, error) {
	file := NewFileBuilder("container.go", params.rootPackageName(), PublicPackage)
	file.AddImport(`"sync"`)
	file.AddImport(params.packageName(InternalPackage))

	var methods bytes.Buffer
	var arguments bytes.Buffer
	var argumentSetters bytes.Buffer
	for _, service := range container.Services {
		parameters := templateParameters{
			ContainerName: container.Name,
			ServiceName:   service.Name,
			ServiceTitle:  service.Title(),
			ServiceType:   service.Type.String(),
		}

		needsImport := false

		if service.IsPublic {
			needsImport = true
			err := publicGetterTemplate.Execute(&methods, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate getter for %s", service.Name)
			}
		}
		if service.HasSetter {
			needsImport = true
			err := publicSetterTemplate.Execute(&methods, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate setter for %s", service.Name)
			}
		}
		if service.IsRequired {
			needsImport = true
			err := argumentTemplate.Execute(&arguments, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate argument for %s", service.Name)
			}
			err = argumentSetterTemplate.Execute(&argumentSetters, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate argument setter for %s", service.Name)
			}
		}

		if needsImport {
			file.AddImport(container.GetImport(service))
		}
	}

	for _, attachedContainer := range container.Containers {
		for _, service := range attachedContainer.Services {
			parameters := templateParameters{
				ContainerName: container.Name,
				ServicePath:   attachedContainer.Title() + "().(" + attachedContainer.Type.String() + ").",
				ServiceName:   service.Name,
				ServiceTitle:  service.Title(),
				ServiceType:   service.Type.String(),
			}

			needsImport := false

			if service.IsPublic {
				needsImport = true
				err := publicGetterTemplate.Execute(&methods, parameters)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to generate getter for %s", service.Name)
				}
			}
			if service.HasSetter {
				needsImport = true
				err := publicSetterTemplate.Execute(&methods, parameters)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to generate setter for %s", service.Name)
				}
			}
			if service.IsRequired {
				needsImport = true
				err := argumentTemplate.Execute(&arguments, parameters)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to generate argument for %s", service.Name)
				}
				err = argumentSetterTemplate.Execute(&argumentSetters, parameters)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to generate argument setter for %s", service.Name)
				}
			}

			if needsImport {
				file.AddImport(container.GetImport(service))
			}
		}
	}

	if arguments.Len() == 0 {
		arguments.WriteString("injectors ...Injector")
	} else {
		arguments.WriteString("\n\tinjectors ...Injector,\n")
		argumentSetters.WriteString("\n")
	}

	err := publicContainerTemplate.Execute(file, containerTemplateParameters{
		ContainerArguments:       arguments.String(),
		ContainerArgumentSetters: argumentSetters.String(),
		ContainerMethods:         methods.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate public container")
	}

	return file.GetFile(), nil
}
