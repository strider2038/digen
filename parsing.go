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

func ParseFile(filename string) (*ContainerDefinition, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse file %s", filename)
	}

	return ParseAST(file)
}

func ParseSource(source string) (*ContainerDefinition, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse source")
	}

	return ParseAST(file)
}

func ParseAST(file *ast.File) (*ContainerDefinition, error) {
	container, err := getContainer(file)
	if err != nil {
		return nil, err
	}

	services, err := createServiceDefinitions(container)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse service definitions")
	}

	imports, err := parseImports(file)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse imports")
	}

	if file.Name == nil {
		return nil, errors.Wrap(ErrParsing, "missing package name")
	}

	definition := &ContainerDefinition{
		Package:  file.Name.Name,
		Imports:  imports,
		Services: services,
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

func createServiceDefinitions(container *ast.StructType) ([]ServiceDefinition, error) {
	services := make([]ServiceDefinition, 0, len(container.Fields.List))

	for _, field := range container.Fields.List {
		fieldType, err := parseFieldType(field)
		if err != nil {
			return nil, err
		}

		service := newServiceDefinition(
			parseFieldName(field),
			fieldType,
			parseFieldTags(field),
		)

		services = append(services, service)
	}

	return services, nil
}

func parseImports(file *ast.File) ([]ImportDefinition, error) {
	imports := make([]ImportDefinition, 0, len(file.Imports))

	for _, spec := range file.Imports {
		imp, err := parseImportDefinition(spec)
		if err != nil {
			return nil, err
		}
		imports = append(imports, imp)
	}

	return imports, nil
}

func parseImportDefinition(spec *ast.ImportSpec) (ImportDefinition, error) {
	definition := ImportDefinition{}

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
			return definition, errors.Wrap(err, "failed to parse import path")
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
	}

	return TypeDefinition{}, errors.Wrap(ErrUnexpectedType, "failed to parse type")
}

func parseFieldTags(field *ast.Field) []string {
	if field.Tag != nil && len(field.Tag.Value) > 1 {
		tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
		return strings.Split(tag.Get("di"), ",")
	}

	return nil
}
