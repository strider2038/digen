package di

import (
	"go/ast"
	"strings"

	"github.com/muonsoft/errors"
)

func ParseFactoryFromSource(source string) (*FactoryFile, error) {
	file, err := parseSource(source)
	if err != nil {
		return nil, err
	}

	return parseFactoryAST(file)
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
