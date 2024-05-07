package di

import (
	"go/ast"
	"slices"
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

	FactoryFileName string // "factory-file" tag
	HasSetter       bool   // "set" tag - will generate setters for internal and public containers
	HasCloser       bool   // "close" tag - generate closer method call
	IsRequired      bool   // "required" tag - will generate argument for public container constructor
	IsPublic        bool   // "public" tag - will generate getter for public container
	IsExternal      bool   // "external" tag - no definition, panic if empty, force public setter
}

func (s ServiceDefinition) Title() string {
	return strings.Title(s.Name)
}

func newServiceDefinition(field *ast.Field, typeDef TypeDefinition) *ServiceDefinition {
	name := parseFieldName(field)
	tags := parseFieldTags(field)

	definition := &ServiceDefinition{
		Name:            name,
		Type:            typeDef,
		FactoryFileName: tags.FactoryFilename,
	}
	if definition.FactoryFileName != "" {
		definition.FactoryFileName += ".go"
	}

	for _, option := range tags.Options {
		switch option {
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

func (d TypeDefinition) ZeroComparison() string {
	if d.IsPointer {
		return " == nil"
	}
	if d.IsBasicType() {
		return d.basicZeroComparison()
	}
	if d.IsTime() {
		return ".IsZero()"
	}
	if d.IsDuration() {
		return " == 0"
	}
	if d.IsURL() {
		return " == url.URL{}"
	}

	return " == nil"
}

func (d TypeDefinition) basicZeroComparison() string {
	switch d.Name {
	case "bool":
		return " == false"
	case "string":
		return " == \"\""
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return " == 0"
	case "float32", "float64":
		return " == 0.0"
	default:
		return " == nil"
	}
}

type FactoryFile struct {
	Imports  map[string]*ImportDefinition
	Services []string
}

type Tags struct {
	Options         []string
	FactoryFilename string
}
