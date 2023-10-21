package di

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"

	"github.com/muonsoft/errors"
)

func ParseDefinitionsFromFile(filename string) (*RootContainerDefinition, error) {
	file, err := parseFile(filename)
	if err != nil {
		return nil, err
	}

	return parseContainerAST(file)
}

func ParseContainerFromSource(source string) (*RootContainerDefinition, error) {
	file, err := parseSource(source)
	if err != nil {
		return nil, err
	}

	return parseContainerAST(file)
}

func ParseFactoryFromSource(source string) (*FactoryFile, error) {
	file, err := parseSource(source)
	if err != nil {
		return nil, err
	}

	return parseFactoryAST(file)
}

func parseFile(filename string) (*ast.File, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, errors.Errorf("parse file %s: %w", filename, err)
	}
	return file, nil
}

func parseSource(source string) (*ast.File, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, errors.Errorf("parse source: %w", err)
	}
	return file, nil
}

func parseContainerAST(file *ast.File) (*RootContainerDefinition, error) {
	container, err := getContainer(file)
	if err != nil {
		return nil, err
	}

	services, containers, err := createFactories(container)
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
	}

	return definition, nil
}

func getContainer(file *ast.File) (*ast.StructType, error) {
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

func createFactories(container *ast.StructType) ([]*ServiceDefinition, []*ContainerDefinition, error) {
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

		if isContainerDefinition(field) {
			internalContainer, err := createContainerDefinition(field)
			if err != nil {
				return nil, nil, err
			}
			containers = append(containers, internalContainer)
		} else {
			service := newServiceDefinition(field, fieldType)
			services = append(services, service)
		}
	}

	return services, containers, nil
}

func isContainerDefinition(field *ast.Field) bool {
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

func createContainerDefinition(field *ast.Field) (*ContainerDefinition, error) {
	fieldType, err := parseFieldType(field)
	if err != nil {
		return nil, err
	}
	fieldType.Package = "internal"

	definition := &ContainerDefinition{
		Name: parseFieldName(field),
		Type: fieldType,
	}

	container, err := parseContainerField(field)
	if err != nil {
		return nil, err
	}

	definition.Services, err = createContainerServiceFactories(container, definition.Name)
	if err != nil {
		return nil, err
	}

	return definition, nil
}

func createContainerServiceFactories(container *ast.StructType, path string) ([]*ServiceDefinition, error) {
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

		if isContainerDefinition(field) {
			return nil, errors.Errorf("%w: %s", ErrNotSupported, "container inside container")
		}

		service := newServiceDefinition(field, fieldType)
		service.Prefix = path
		services = append(services, service)
	}

	return services, nil
}

func parseImports(file *ast.File) (map[string]*ImportDefinition, error) {
	imports := make(map[string]*ImportDefinition, len(file.Imports))

	for _, spec := range file.Imports {
		imp, err := parseImportDefinition(spec)
		if err != nil {
			return nil, err
		}
		imports[imp.ID] = imp
	}

	return imports, nil
}

func parseImportDefinition(spec *ast.ImportSpec) (*ImportDefinition, error) {
	definition := &ImportDefinition{}

	if spec.Name != nil {
		definition.Name = spec.Name.Name
	}

	if spec.Path != nil {
		definition.Path = spec.Path.Value
	}

	if definition.Name != "" {
		definition.ID = definition.Name
	} else {
		path, err := strconv.Unquote(definition.Path)
		if err != nil {
			return nil, errors.Errorf("parse import path: %w", err)
		}
		elements := strings.Split(path, "/")
		if len(elements) > 0 {
			definition.ID = elements[len(elements)-1]
		}
	}

	return definition, nil
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

	case *ast.Ident:
		return TypeDefinition{Name: t.Name}, nil
	}

	return TypeDefinition{}, errors.Errorf("%w: %s", ErrUnexpectedType, "parse type")
}

func parseFieldTags(field *ast.Field) Tags {
	if field.Tag == nil || len(field.Tag.Value) == 0 {
		return Tags{}
	}

	tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])

	return Tags{
		Options:         strings.Split(tag.Get("di"), ","),
		FactoryFilename: tag.Get("factory-file"),
	}
}

func parseContainerField(field *ast.Field) (*ast.StructType, error) {
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

func parseFactoryAST(file *ast.File) (*FactoryFile, error) {
	imports, err := parseImports(file)
	if err != nil {
		return nil, errors.Errorf("parse imports: %w", err)
	}

	var services []string

	for name, object := range file.Scope.Objects {
		if object.Kind == ast.Fun && strings.HasPrefix(name, "Create") {
			services = append(services, strings.TrimPrefix(name, "Create"))
		}
	}

	factory := &FactoryFile{
		Imports:  imports,
		Services: services,
	}

	return factory, nil
}

func validateInternalContainer(container *ast.StructType) error {
	if len(container.Fields.List) == 0 {
		return errors.Errorf("%w: %s", ErrInvalidDefinition, "container must not be empty")
	}

	return nil
}
