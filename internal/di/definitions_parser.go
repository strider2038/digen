package di

import (
	"go/ast"
	"reflect"
	"strings"

	"github.com/muonsoft/errors"
)

func ParseDefinitionsFromFile(filename string) (*RootContainerDefinition, error) {
	parser := &DefinitionsParser{}

	return parser.ParseFile(filename)
}

func ParseContainerFromSource(source string) (*RootContainerDefinition, error) {
	parser := &DefinitionsParser{}

	return parser.ParseSource(source)
}

type DefinitionsParser struct {
	lastID int
}

func (p *DefinitionsParser) ParseFile(filename string) (*RootContainerDefinition, error) {
	file, err := parseFile(filename)
	if err != nil {
		return nil, err
	}

	return p.parseContainerAST(file)
}

func (p *DefinitionsParser) ParseSource(source string) (*RootContainerDefinition, error) {
	file, err := parseSource(source)
	if err != nil {
		return nil, err
	}

	return p.parseContainerAST(file)
}

func (p *DefinitionsParser) parseContainerAST(file *ast.File) (*RootContainerDefinition, error) {
	container, err := p.getContainer(file)
	if err != nil {
		return nil, err
	}

	services, containers, err := p.parseDefinitions(container)
	if err != nil {
		return nil, errors.Errorf("parse definitions: %w", err)
	}

	imports, err := parseImports(file)
	if err != nil {
		return nil, errors.Errorf("parse imports: %w")
	}

	if file.Name == nil {
		return nil, errors.Errorf("%w: %s", ErrParsing, "missing package name")
	}

	definition := &RootContainerDefinition{
		Name:       "Container",
		Package:    file.Name.Name,
		Imports:    imports,
		Services:   services,
		Containers: containers,
		Factories:  make(map[string]*FactoryDefinition, 0),
	}

	return definition, nil
}

func (p *DefinitionsParser) getContainer(file *ast.File) (*ast.StructType, error) {
	container := file.Scope.Objects["Container"]
	if container == nil {
		return nil, errors.Wrap(ErrContainerNotFound)
	}

	containerType, ok := container.Decl.(*ast.TypeSpec)
	if !ok {
		return nil, errors.Wrap(ErrUnexpectedType)
	}
	containerStruct, ok := containerType.Type.(*ast.StructType)
	if !ok {
		return nil, errors.Wrap(ErrUnexpectedType)
	}

	return containerStruct, nil
}

func (p *DefinitionsParser) parseDefinitions(container *ast.StructType) ([]*ServiceDefinition, []*ContainerDefinition, error) {
	services := make([]*ServiceDefinition, 0)
	containers := make([]*ContainerDefinition, 0)

	for _, field := range container.Fields.List {
		fieldType, err := parseFieldType(field)
		if err != nil {
			return nil, nil, err
		}
		if fieldType.Name == "error" {
			continue
		}

		if p.isContainerDefinition(field) {
			internalContainer, err := p.createContainerDefinition(field)
			if err != nil {
				return nil, nil, err
			}
			containers = append(containers, internalContainer)
		} else {
			service := p.createServiceDefinition(field, fieldType)
			services = append(services, service)
		}
	}

	return services, containers, nil
}

func (p *DefinitionsParser) isContainerDefinition(field *ast.Field) bool {
	if id, ok := field.Type.(*ast.Ident); ok {
		if id.Obj != nil {
			if t, ok := id.Obj.Decl.(*ast.TypeSpec); ok {
				_, ok := t.Type.(*ast.StructType)

				return ok
			}
		}
	}

	return false
}

func (p *DefinitionsParser) createContainerDefinition(field *ast.Field) (*ContainerDefinition, error) {
	fieldType, err := parseFieldType(field)
	if err != nil {
		return nil, err
	}
	fieldType.Package = "internal"

	definition := &ContainerDefinition{
		Name: parseFieldName(field),
		Type: fieldType,
	}

	container, err := p.parseContainerField(field)
	if err != nil {
		return nil, err
	}

	definition.Services, err = p.parseServiceDefinitions(container, definition.Name)
	if err != nil {
		return nil, err
	}

	return definition, nil
}

func (p *DefinitionsParser) createServiceDefinition(field *ast.Field, typeDef TypeDefinition) *ServiceDefinition {
	name := parseFieldName(field)
	tags := parseFieldTags(field)

	definition := &ServiceDefinition{
		ID:              p.nextID(),
		Name:            name,
		Type:            typeDef,
		FactoryFileName: tags.FactoryFilename,
		PublicName:      tags.PublicName,
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

func (p *DefinitionsParser) parseServiceDefinitions(container *ast.StructType, path string) ([]*ServiceDefinition, error) {
	services := make([]*ServiceDefinition, 0)
	err := validateInternalContainer(container)
	if err != nil {
		return nil, err
	}

	for _, field := range container.Fields.List {
		fieldType, err := parseFieldType(field)
		if err != nil {
			return nil, err
		}

		if p.isContainerDefinition(field) {
			return nil, errors.Errorf("%w: %s", ErrNotSupported, "container inside container")
		}

		service := p.createServiceDefinition(field, fieldType)
		service.Prefix = path
		services = append(services, service)
	}

	return services, nil
}

func (p *DefinitionsParser) parseContainerField(field *ast.Field) (*ast.StructType, error) {
	fieldDeclaration, ok := field.Names[0].Obj.Decl.(*ast.Field)
	if !ok {
		return nil, errors.Errorf("%w: %s", ErrParsing, "unexpected field declaration")
	}
	containerType, ok := fieldDeclaration.Type.(*ast.Ident)
	if !ok {
		return nil, errors.Errorf("%w: %s", ErrParsing, "unexpected container declaration")
	}
	typeSpecification, ok := containerType.Obj.Decl.(*ast.TypeSpec)
	if !ok {
		return nil, errors.Errorf("%w: %s", ErrParsing, "unexpected container type specification")
	}
	container, ok := typeSpecification.Type.(*ast.StructType)
	if !ok {
		return nil, errors.Errorf("%w: %s", ErrParsing, "container type must be struct")
	}

	return container, nil
}

func (p *DefinitionsParser) nextID() int {
	id := p.lastID
	p.lastID++

	return id
}

func parseFieldName(field *ast.Field) string {
	var s strings.Builder
	for _, ident := range field.Names {
		s.WriteString(ident.Name)
	}

	return s.String()
}

func parseFieldType(field *ast.Field) (TypeDefinition, error) {
	return parseTypeDefinition(field.Type)
}

func parseTypeDefinition(expr ast.Expr) (TypeDefinition, error) {
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		definition := TypeDefinition{}
		ident, ok := t.X.(*ast.Ident)
		if !ok {
			return definition, errors.Errorf("%w: %s", ErrUnexpectedType, "parse package")
		}
		definition.Package = ident.Name
		definition.Name = t.Sel.Name

		return definition, nil

	case *ast.StarExpr:
		definition, err := parseTypeDefinition(t.X)
		if err != nil {
			return definition, err
		}
		if definition.IsPointer {
			return definition, errors.Errorf("%w: %s", ErrNotSupported, "double pointers")
		}
		definition.IsPointer = true

		return definition, nil

	case *ast.ArrayType:
		definition, err := parseTypeDefinition(t.Elt)
		if err != nil {
			return definition, err
		}
		if t.Len != nil {
			return definition, errors.Errorf("%w: %s", ErrNotSupported, "array with length")
		}
		definition.IsSlice = true

		return definition, nil

	case *ast.MapType:
		definition, err := parseTypeDefinition(t.Value)
		if err != nil {
			return definition, err
		}
		key, err := parseTypeDefinition(t.Key)
		if err != nil {
			return definition, err
		}
		definition.Key = &key

		return definition, nil

	case *ast.Ident:
		return TypeDefinition{Name: t.Name}, nil
	}

	return TypeDefinition{}, errors.Errorf("%w: %s", ErrUnexpectedType, "parse type")
}

type Tags struct {
	Options         []string
	FactoryFilename string
	PublicName      string
}

func parseFieldTags(field *ast.Field) Tags {
	if field.Tag == nil || len(field.Tag.Value) == 0 {
		return Tags{}
	}

	tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])

	return Tags{
		Options:         strings.Split(tag.Get("di"), ","),
		FactoryFilename: tag.Get("factory_file"),
		PublicName:      tag.Get("public_name"),
	}
}

func validateInternalContainer(container *ast.StructType) error {
	if len(container.Fields.List) == 0 {
		return errors.Errorf("%w: %s", ErrInvalidDefinition, "container must not be empty")
	}

	return nil
}

type FuncDeclaration struct {
	ReturnsErr bool
}

func parseFuncDeclaration(decl *ast.FuncDecl) (FuncDeclaration, error) {
	declaration := FuncDeclaration{}

	if decl.Type.Results != nil {
		for _, field := range decl.Type.Results.List {
			if id, ok := field.Type.(*ast.Ident); ok && id.Name == "error" {
				declaration.ReturnsErr = true
			}
		}
	}

	return declaration, nil
}
