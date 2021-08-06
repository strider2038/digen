package digen

import (
	"bytes"
	"strings"

	"github.com/pkg/errors"
)

type PackageType int

const (
	unknownPackage PackageType = iota
	PublicPackage
	InternalPackage
	DefinitionsPackage
	lastPackage
)

type File struct {
	Package PackageType
	Name    string
	Content []byte
}

func (f *File) WithHeading(heading []byte) *File {
	return &File{
		Package: f.Package,
		Name:    f.Name,
		Content: append(heading, f.Content...),
	}
}

type GenerationParameters struct {
}

func DefaultGenerationParameters() GenerationParameters {
	return GenerationParameters{}
}

func Generate(container *ContainerDefinition, params GenerationParameters) ([]*File, error) {
	files := make([]*File, 0)

	for _, generate := range generators {
		file, err := generate(container, params)
		if errors.Is(err, errFileIgnored) {
			continue
		}
		if err != nil {
			return nil, errors.WithMessage(err, "failed to generate getters file")
		}

		files = append(files, file)
	}

	return files, nil
}

func GenerateContainer(params GenerationParameters) (*File, error) {
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

var generators = [...]func(container *ContainerDefinition, params GenerationParameters) (*File, error){
	generateContainerGettersFile,
	generateContainerSettersFile,
	generateContainerCloseFile,
	generateDefinitionsContainerFile,
	generatePublicContainerFile,
}

func generateContainerGettersFile(container *ContainerDefinition, params GenerationParameters) (*File, error) {
	var buffer bytes.Buffer

	buffer.WriteString("package " + container.Package + "\n\n")
	buffer.WriteString("import (\n")
	for _, imp := range container.Imports {
		buffer.WriteString("\t" + imp.String() + "\n")
	}
	// todo: import definitions package
	buffer.WriteString(")\n")

	for _, service := range container.Services {
		err := getterTemplate.Execute(&buffer, templateParameters{
			ContainerName: container.Name,
			ServiceName:   service.Name,
			ServiceTitle:  service.Title(),
			ServiceType:   service.Type.String(),
			HasDefinition: !service.IsRequired,
			PanicOnNil:    service.IsExternal,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate getter for %s", service.Name)
		}
	}

	return &File{
		Package: InternalPackage,
		Name:    "container_get.go",
		Content: buffer.Bytes(),
	}, nil
}

func generateContainerSettersFile(container *ContainerDefinition, params GenerationParameters) (*File, error) {
	var buffer bytes.Buffer

	buffer.WriteString("package " + container.Package + "\n\n")
	buffer.WriteString("import (\n")
	for _, imp := range container.Imports {
		buffer.WriteString("\t" + imp.String() + "\n")
	}
	// todo: import definitions package
	buffer.WriteString(")\n")

	count := 0

	for _, service := range container.Services {
		if service.HasSetter || service.IsExternal || service.IsRequired {
			count++
			err := setterTemplate.Execute(&buffer, templateParameters{
				ContainerName: container.Name,
				ServiceName:   service.Name,
				ServiceTitle:  service.Title(),
				ServiceType:   service.Type.String(),
			})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate setter for %s", service.Name)
			}
		}
	}

	if count == 0 {
		return nil, errFileIgnored
	}

	return &File{
		Package: InternalPackage,
		Name:    "container_set.go",
		Content: buffer.Bytes(),
	}, nil
}

func generateContainerCloseFile(container *ContainerDefinition, params GenerationParameters) (*File, error) {
	var buffer bytes.Buffer

	buffer.WriteString("package " + container.Package + "\n\n")

	buffer.WriteString("func (c *Container) Close() {")
	for _, service := range container.Services {
		if service.HasCloser {
			err := closerTemplate.Execute(&buffer, templateParameters{ServiceName: service.Name})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate closer for %s", service.Name)
			}
		}
	}
	buffer.WriteString("}\n")

	return &File{
		Package: InternalPackage,
		Name:    "container_close.go",
		Content: buffer.Bytes(),
	}, nil
}

func generateDefinitionsContainerFile(container *ContainerDefinition, params GenerationParameters) (*File, error) {
	var buffer bytes.Buffer

	buffer.WriteString("package definitions\n\n")
	buffer.WriteString("import (\n")
	for _, imp := range container.Imports {
		buffer.WriteString("\t" + imp.String() + "\n")
	}
	buffer.WriteString(")\n\n")

	buffer.WriteString("type Container interface {\n")
	buffer.WriteString("\tSetError(err error)\n\n")
	for _, service := range container.Services {
		buffer.WriteString("\t" + strings.Title(service.Name) + "() " + service.Type.String() + "\n")
	}
	buffer.WriteString("}\n")

	return &File{
		Package: DefinitionsPackage,
		Name:    "container.go",
		Content: buffer.Bytes(),
	}, nil
}

func generatePublicContainerFile(container *ContainerDefinition, params GenerationParameters) (*File, error) {
	var buffer bytes.Buffer

	buffer.WriteString("package di\n\n")
	buffer.WriteString("import (\n")
	buffer.WriteString("\t\"sync\"\n\n")
	for _, imp := range container.Imports {
		buffer.WriteString("\t" + imp.String() + "\n")
	}
	buffer.WriteString(")\n")

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

		if service.IsPublic {
			err := publicGetterTemplate.Execute(&methods, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate getter for %s", service.Name)
			}
		}
		if service.HasSetter {
			err := publicSetterTemplate.Execute(&methods, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate setter for %s", service.Name)
			}
		}
		if service.IsRequired {
			err := argumentTemplate.Execute(&arguments, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate argument for %s", service.Name)
			}
			err = argumentSetterTemplate.Execute(&argumentSetters, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to generate argument setter for %s", service.Name)
			}
		}
	}

	if arguments.Len() == 0 {
		arguments.WriteString("injectors ...Injector")
	} else {
		arguments.WriteString("\n\tinjectors ...Injector,\n")
		argumentSetters.WriteString("\n")
	}

	err := publicContainerTemplate.Execute(&buffer, containerTemplateParameters{
		ContainerArguments:       arguments.String(),
		ContainerArgumentSetters: argumentSetters.String(),
		ContainerMethods:         methods.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate public container")
	}

	return &File{
		Package: PublicPackage,
		Name:    "container.go",
		Content: buffer.Bytes(),
	}, nil
}
