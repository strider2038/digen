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
}

func (c RootContainerDefinition) Type(definition TypeDefinition) func(statement *jen.Statement) {
	return func(statement *jen.Statement) {
		packageName := c.PackageName(definition)
		if definition.IsPointer {
			statement.Op("*").Qual(packageName, definition.Name)
		} else {
			statement.Qual(packageName, definition.Name)
		}
	}
}

func (c RootContainerDefinition) PackageName(definition TypeDefinition) string {
	packageName := ""
	if imp, ok := c.Imports[definition.Package]; ok {
		// todo: check parsing to remove cutset
		packageName = strings.Trim(imp.Path, `"`)
	}

	return packageName
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
	IsExternal bool // "external" tag - no definition, panic if empty, force public setter
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

func (d TypeDefinition) ZeroComparison() func(statement *jen.Statement) {
	return func(statement *jen.Statement) {
		switch {
		case d.IsPointer:
			statement.Op("==").Nil()
		case d.IsBasicType():
			d.basicZeroComparison(statement)
		case d.IsTime():
			statement.Dot("IsZero").Call()
		case d.IsDuration():
			statement.Op("==").Lit(0)
		case d.IsURL():
			statement.Op("==").New(jen.Qual("net/url", "URL"))
		default:
			statement.Op("==").Nil()
		}
	}
}

func (d TypeDefinition) basicZeroComparison(statement *jen.Statement) {
	switch d.Name {
	case "bool":
		statement.Op("==").False()
	case "string":
		statement.Op("==").Lit("")
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		statement.Op("==").Lit(0)
	case "float32", "float64":
		statement.Op("==").Lit(0.0)
	default:
		statement.Op("==").Nil()
	}
}

type FactoryFile struct {
	Imports  map[string]*ImportDefinition
	Services []string
}

type Tags struct {
	Options         []string
	FactoryFilename string
	PublicName      string
}
