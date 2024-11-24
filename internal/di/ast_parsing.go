package di

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"github.com/muonsoft/errors"
)

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
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return nil, errors.Errorf("parse import path: %w", err)
		}

		definition.Path = path
	}

	if definition.Name != "" {
		definition.ID = definition.Name
	} else {
		elements := strings.Split(definition.Path, "/")
		if len(elements) > 0 {
			definition.ID = elements[len(elements)-1]
		}
	}

	return definition, nil
}
