package digen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func ParseContainerFromFile(filename string) (*RootContainerDefinition, error) {
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

func ParseFactoryFromFile(filename string) (*FactoryFile, error) {
	file, err := parseFile(filename)
	if err != nil {
		return nil, err
	}

	return parseFactoryAST(file)
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
		return nil, errors.Wrapf(err, "failed to parse file %s", filename)
	}
	return file, nil
}

func parseSource(source string) (*ast.File, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse source")
	}
	return file, nil
}

func parseContainerAST(file *ast.File) (*RootContainerDefinition, error) {
	container, err := getContainer(file)
	if err != nil {
		return nil, err
	}

	services, containers, err := createDefinitions(container)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse definitions")
	}

	imports, err := parseImports(file)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse imports")
	}

	if file.Name == nil {
		return nil, errors.Wrap(ErrParsing, "missing package name")
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
		return nil, errors.WithStack(ErrContainerNotFound)
	}

	containerType, ok := container.Decl.(*ast.TypeSpec)
	if !ok {
		return nil, errors.WithStack(ErrUnexpectedType)
	}
	containerStruct, ok := containerType.Type.(*ast.StructType)
	if !ok {
		return nil, errors.WithStack(ErrUnexpectedType)
	}

	return containerStruct, nil
}

func createDefinitions(container *ast.StructType) ([]*ServiceDefinition, []*ContainerDefinition, error) {
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

		tags := parseFieldTags(field)
		if tags.Contains("container") {
			internalContainer, err := createContainerDefinition(field)
			if err != nil {
				return nil, nil, err
			}
			containers = append(containers, internalContainer)
		} else {
			service := newServiceDefinition(parseFieldName(field), fieldType, tags)
			services = append(services, service)
		}
	}

	return services, containers, nil
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

	definition.Services, err = createContainerServiceDefinitions(container)
	if err != nil {
		return nil, err
	}

	return definition, nil
}

func createContainerServiceDefinitions(container *ast.StructType) ([]*ServiceDefinition, error) {
	services := make([]*ServiceDefinition, 0)
	err := validateInternalContainer(container)
	if err != nil {
		return nil, err
	}

	for i, field := range container.Fields.List {
		if i == 0 {
			continue
		}

		fieldType, err := parseFieldType(field)
		if err != nil {
			return nil, err
		}

		tags := parseFieldTags(field)
		if tags.Contains("container") {
			return nil, errors.Wrap(ErrNotSupported, "container inside container")
		}

		service := newServiceDefinition(parseFieldName(field), fieldType, tags)
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
			return nil, errors.Wrap(err, "failed to parse import path")
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
			return definition, errors.Wrap(ErrUnexpectedType, "failed to parse package")
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
			return definition, errors.Wrap(ErrNotSupported, "double pointers")
		}
		definition.IsPointer = true

		return definition, nil

	case *ast.Ident:
		return TypeDefinition{Name: t.Name}, nil
	}

	return TypeDefinition{}, errors.Wrap(ErrUnexpectedType, "failed to parse type")
}

func parseFieldTags(field *ast.Field) Tags {
	if field.Tag != nil && len(field.Tag.Value) > 1 {
		tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])

		return strings.Split(tag.Get("di"), ",")
	}

	return nil
}

func parseContainerField(field *ast.Field) (*ast.StructType, error) {
	fieldDeclaration, ok := field.Names[0].Obj.Decl.(*ast.Field)
	if !ok {
		return nil, errors.Wrap(ErrParsing, "unexpected field declaration")
	}
	containerPointer, ok := fieldDeclaration.Type.(*ast.StarExpr)
	if !ok {
		return nil, errors.Wrap(ErrParsing, "container type must be pointer")
	}
	containerType, ok := containerPointer.X.(*ast.Ident)
	if !ok {
		return nil, errors.Wrap(ErrParsing, "unexpected container declaration")
	}
	typeSpecification, ok := containerType.Obj.Decl.(*ast.TypeSpec)
	if !ok {
		return nil, errors.Wrap(ErrParsing, "unexpected container type specification")
	}
	container, ok := typeSpecification.Type.(*ast.StructType)
	if !ok {
		return nil, errors.Wrap(ErrParsing, "container type must be struct")
	}

	return container, nil
}

func parseFactoryAST(file *ast.File) (*FactoryFile, error) {
	imports, err := parseImports(file)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse imports")
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
		return errors.Wrap(ErrInvalidDefinition, "container must not be empty")
	}
	field := container.Fields.List[0]

	containerPointer, ok := field.Type.(*ast.StarExpr)
	if !ok {
		return errors.Wrap(ErrInvalidDefinition, "internal container must embed root container as a pointer in the first field")
	}
	rootContainer, ok := containerPointer.X.(*ast.Ident)
	if !ok {
		return errors.Wrap(ErrInvalidDefinition, "internal container must embed root container as a pointer in the first field")
	}
	if rootContainer.Name != "Container" {
		return errors.Wrap(ErrInvalidDefinition, "internal container must embed root container as a pointer in the first field")
	}
	if len(field.Names) > 0 {
		return errors.Wrap(ErrInvalidDefinition, "internal container must embed root container as a pointer in the first field")
	}

	return nil
}
