package di

import (
	"slices"
	"strings"

	"github.com/dave/jennifer/jen"
)

type RootContainerDefinition struct {
	Name       string
	Package    string
	Imports    map[string]*ImportDefinition
	Services   []*ServiceDefinition
	Containers []*ContainerDefinition
	Factories  map[string]*FactoryDefinition
}

func (c RootContainerDefinition) Type(definition TypeDefinition) func(statement *jen.Statement) {
	return func(statement *jen.Statement) {
		packageName := c.PackageName(definition)
		if definition.IsPointer {
			statement = statement.Op("*")
		} else if definition.IsSlice {
			statement = statement.Op("[]")
		} else if definition.IsMap() {
			statement = statement.Map(jen.Id(definition.Key.Name))
		}
		statement.Qual(packageName, definition.Name)
	}
}

func (c RootContainerDefinition) PackageName(definition TypeDefinition) string {
	if imp, ok := c.Imports[definition.Package]; ok {
		return imp.Path
	}

	return ""
}

func (c RootContainerDefinition) ServicesCount() int {
	count := len(c.Services)

	for _, container := range c.Containers {
		count += len(container.Services)
	}

	return count
}

type ImportDefinition struct {
	ID   string
	Name string
	Path string
}

func (d ImportDefinition) String() string {
	if d.Name != "" {
		return d.Name + " " + d.Path
	}

	return d.Path
}

type ServiceDefinition struct {
	ID     int
	Prefix string
	Name   string
	Type   TypeDefinition

	PublicName      string // "public_name" tag
	FactoryFileName string // "factory_file" tag

	// options from tag "di"
	HasSetter  bool // "set" tag - will generate setters for internal and public containers
	HasCloser  bool // "close" tag - generate closer method call
	IsRequired bool // "required" tag - will generate argument for public container constructor
	IsPublic   bool // "public" tag - will generate getter for public container
}

func (s ServiceDefinition) Title() string {
	return strings.Title(s.Name)
}

func (s ServiceDefinition) PublicTitle() string {
	if s.PublicName != "" {
		return strings.Title(s.PublicName)
	}

	return s.Title()
}

type ContainerDefinition struct {
	Name     string
	Type     TypeDefinition
	Services []*ServiceDefinition
}

func (c ContainerDefinition) Title() string {
	return strings.Title(c.Name)
}

type TypeDefinition struct {
	IsPointer bool
	IsSlice   bool
	Package   string
	Name      string
	Key       *TypeDefinition
}

func (d TypeDefinition) IsMap() bool {
	return d.Key != nil
}

var basicTypes = []string{
	"string",
	"int",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"int8",
	"int16",
	"int32",
	"int64",
	"float32",
	"float64",
	"bool",
}

func (d TypeDefinition) IsBasicType() bool {
	return d.Package == "" && slices.Contains(basicTypes, d.Name)
}

func (d TypeDefinition) IsTime() bool {
	return d.Package == "time" && d.Name == "Time"
}

func (d TypeDefinition) IsDuration() bool {
	return d.Package == "time" && d.Name == "Duration"
}

func (d TypeDefinition) IsURL() bool {
	return d.Package == "url" && d.Name == "URL"
}

func (d TypeDefinition) String() string {
	var s strings.Builder

	if d.IsPointer {
		s.WriteString("*")
	}
	s.WriteString(d.Package)
	if d.Package != "" {
		s.WriteString(".")
	}
	s.WriteString(d.Name)

	return s.String()
}

type FactoryDefinitions struct {
	Imports   map[string]*ImportDefinition
	Factories map[string]*FactoryDefinition
}

func NewFactoryDefinitions() *FactoryDefinitions {
	return &FactoryDefinitions{
		Imports:   map[string]*ImportDefinition{},
		Factories: map[string]*FactoryDefinition{},
	}
}

func (d *FactoryDefinitions) merge(df *FactoryDefinitions) {
	for k, v := range df.Factories {
		d.Factories[k] = v
	}
	for k, v := range df.Imports {
		d.Imports[k] = v
	}
}

type FactoryDefinition struct {
	Name         string
	ReturnsError bool
}
