package di

import (
	"go/ast"
	iofs "io/fs"
	"strings"

	"github.com/muonsoft/errors"
	"github.com/spf13/afero"
)

func parseFactoriesFromDir(fs afero.Fs, dir string) (*FactoryDefinitions, error) {
	definitions := NewFactoryDefinitions()
	if !isFileExist(fs, dir) {
		return definitions, nil
	}

	err := afero.Walk(fs, dir, func(path string, d iofs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, err := parseFile(fs, path)
		if err != nil {
			return err
		}
		df, err := parseFactoriesAST(file)
		if err != nil {
			return err
		}
		definitions.merge(df)

		return nil
	})
	if err != nil {
		return nil, errors.Errorf("walk dir %q: %w", dir, err)
	}

	return definitions, nil
}

func ParseFactoriesFromSource(source string) (*FactoryDefinitions, error) {
	file, err := parseSource(source)
	if err != nil {
		return nil, err
	}

	return parseFactoriesAST(file)
}

func parseFactoriesAST(file *ast.File) (*FactoryDefinitions, error) {
	imports, err := parseImports(file)
	if err != nil {
		return nil, errors.Errorf("parse imports: %w", err)
	}

	factories := make(map[string]*FactoryDefinition, len(file.Scope.Objects))

	for name, object := range file.Scope.Objects {
		if funcDecl, ok := object.Decl.(*ast.FuncDecl); ok && object.Kind == ast.Fun && strings.HasPrefix(name, "Create") {
			f, err := parseFuncDeclaration(funcDecl)
			if err != nil {
				return nil, errors.Errorf("parse func declaration: %w", err)
			}
			factoryName := strings.TrimPrefix(name, "Create")
			factories[factoryName] = &FactoryDefinition{
				Name:         factoryName,
				ReturnsError: f.ReturnsErr,
			}
		}
	}

	return &FactoryDefinitions{
		Imports:   imports,
		Factories: factories,
	}, nil
}
