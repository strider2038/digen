package digen

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

type PackageType int

const (
	UnknownPackage PackageType = iota
	PublicPackage
	InternalPackage
	DefinitionsPackage
)

type File struct {
	Package PackageType
	Name    string
	Content []byte
}

func (f *File) WriteTo(dir string) error {
	err := os.WriteFile(dir+"/"+f.Name, f.Content, 0644)
	if err != nil {
		return errors.WithMessagef(err, "failed to write %s", f.Name)
	}

	return nil
}

type GenerationParameters struct {
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

var generators = [...]func(container *ContainerDefinition, params GenerationParameters) (*File, error){
	generateContainerGettersFile,
	generateContainerSettersFile,
	generateContainerCloseFile,
	generateDefinitionsContainerFile,
	// definitions
	// public container
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
			ContainerName: "Container", // todo: replace by real name
			ServiceName:   service.Name,
			ServiceTitle:  service.Title(),
			ServiceType:   service.Type.String(),
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
		if !service.HasSetter {
			continue
		}
		count++
		err := setterTemplate.Execute(&buffer, templateParameters{
			ContainerName: "Container", // todo: replace by real name
			ServiceName:   service.Name,
			ServiceTitle:  service.Title(),
			ServiceType:   service.Type.String(),
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate setter for %s", service.Name)
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
	buffer.WriteString("import (\n")
	for _, imp := range container.Imports {
		buffer.WriteString("\t" + imp.String() + "\n")
	}
	buffer.WriteString(")\n\n")

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
	for _, service := range container.Services {
		buffer.WriteString("\t" + strings.Title(service.Name) + "() " + service.Type.String() + "\n")
	}
	buffer.WriteString("}\n")

	return &File{
		Package: DefinitionsPackage,
		Name:    "definitions/container.go",
		Content: buffer.Bytes(),
	}, nil
}

type templateParameters struct {
	ContainerName string
	ServiceName   string
	ServiceTitle  string
	ServiceType   string
}

var getterTemplate = template.Must(template.New("getter").Parse(`
func (c *{{.ContainerName}}) {{.ServiceTitle}}() {{.ServiceType}} {
	if c.{{.ServiceName}} == nil {
		c.{{.ServiceName}} = definitions.Create{{.ServiceTitle}}(c)
	}
	return c.{{.ServiceName}}
}
`))

var setterTemplate = template.Must(template.New("setter").Parse(`
func (c *{{.ContainerName}}) Set{{.ServiceTitle}}(s {{.ServiceType}}) {
	c.{{.ServiceName}} = s
}
`))

var closerTemplate = template.Must(template.New("closer").Parse(`
	if c.{{.ServiceName}} != nil {
		c.{{.ServiceName}}.Close()
	}
`))
