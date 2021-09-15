package di

import (
	"go/ast"
	"strings"
)

type RootContainerDefinition struct {
	Name       string
	Package    string
	Imports    map[string]*ImportDefinition
	Services   []*ServiceDefinition
	Containers []*ContainerDefinition
}

func (c RootContainerDefinition) GetImport(s *ServiceDefinition) string {
	imp := c.Imports[s.Type.Package]
	if imp == nil {
		return ""
	}
	return imp.String()
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
	Prefix string
	Name   string
	Type   TypeDefinition

	HasSetter  bool // "set" tag - will generate setters for internal and public containers
	HasCloser  bool // "close" tag - generate closer method call
	IsRequired bool // "required" tag - will generate argument for public container constructor
	IsPublic   bool // "public" tag - will generate getter for public container
	IsExternal bool // "external" tag - no definition, panic if empty, force public setter
}

func (s ServiceDefinition) Title() string {
	return strings.Title(s.Name)
}

func newServiceDefinition(field *ast.Field, typeDef TypeDefinition) *ServiceDefinition {
	name := parseFieldName(field)
	tags := parseFieldTags(field)

	definition := &ServiceDefinition{Name: name, Type: typeDef}

	for _, tag := range tags {
		switch tag {
		case "set":
			definition.HasSetter = true
		case "close":
			definition.HasCloser = true
		case "required":
			definition.IsRequired = true
		case "public":
			definition.IsPublic = true
		case "external":
			definition.IsExternal = true
		}
	}

	return definition
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
	Package   string
	Name      string
}

func (d TypeDefinition) String() string {
	var s strings.Builder

	if d.IsPointer {
		s.WriteString("*")
	}
	s.WriteString(d.Package)
	s.WriteString(".")
	s.WriteString(d.Name)

	return s.String()
}

type FactoryFile struct {
	Imports  map[string]*ImportDefinition
	Services []string
}

type Tags []string

func (tags Tags) Contains(tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}

	return false
}
