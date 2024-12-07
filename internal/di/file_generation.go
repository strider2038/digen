package di

import (
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/muonsoft/errors"
)

func GenerateDefinitionsContainerFile() *File {
	return &File{
		Package: DefinitionsPackage,
		Name:    "container.go",
		Content: []byte(definitionsContainerFileSkeleton),
	}
}

type GenerationParameters struct {
	RootPackage   string
	ErrorHandling ErrorHandling
}

type ErrorHandling struct {
	New  ErrorOptions
	Join ErrorOptions
	Wrap ErrorOptions
}

type ErrorOptions struct {
	Package  string
	Function string
	Verb     string
}

func (w ErrorHandling) Defaults() ErrorHandling {
	if w.New.Package == "" {
		w.New.Package = "fmt"
	}
	if w.New.Function == "" {
		w.New.Function = "Errorf"
	}
	if w.Join.Package == "" {
		w.Join.Package = "errors"
	}
	if w.Join.Function == "" {
		w.Join.Function = "Join"
	}
	if w.Wrap.Package == "" {
		w.Wrap.Package = "fmt"
	}
	if w.Wrap.Function == "" {
		w.Wrap.Function = "Errorf"
	}
	if w.Wrap.Verb == "" {
		w.Wrap.Verb = "%w"
	}

	return w
}

func (params *GenerationParameters) rootPackageName() string {
	path := strings.Split(params.RootPackage, "/")
	if len(path) == 0 {
		return ""
	}
	return path[len(path)-1]
}

func (params *GenerationParameters) packageName(packageType PackageType) string {
	return strings.Trim(strconv.Quote(params.RootPackage+"/"+packageDirs[packageType]), `"`)
}

func (params *GenerationParameters) wrapError(message string, errorIdentifier jen.Code) *jen.Statement {
	path := params.ErrorHandling.Wrap.Package
	funcName := params.ErrorHandling.Wrap.Function
	verb := params.ErrorHandling.Wrap.Verb

	return jen.Qual(path, funcName).Call(jen.Lit(message+": "+verb), errorIdentifier)
}

func (params *GenerationParameters) joinErrors(errs ...jen.Code) *jen.Statement {
	path := params.ErrorHandling.Join.Package
	funcName := params.ErrorHandling.Join.Function

	return jen.Qual(path, funcName).Call(errs...)
}

func GenerateFiles(container *RootContainerDefinition, params GenerationParameters) ([]*File, error) {
	return NewFileGenerator(container, params).GenerateFiles()
}

type FileGenerator struct {
	container *RootContainerDefinition
	params    GenerationParameters
}

func NewFileGenerator(
	container *RootContainerDefinition,
	params GenerationParameters,
) *FileGenerator {
	return &FileGenerator{
		container: container,
		params:    params,
	}
}

func (g *FileGenerator) GenerateFiles() ([]*File, error) {
	files := make([]*File, 0)

	generators := [...]func() (*File, error){
		NewInternalContainerGenerator(g.container, g.params).Generate,
		NewLookupContainerGenerator(g.container).Generate,
		NewPublicContainerGenerator(g.container, g.params).Generate,
	}

	for _, generate := range generators {
		file, err := generate()
		if err != nil {
			return nil, errors.Errorf("generate file: %w", err)
		}

		files = append(files, file)
	}

	return files, nil
}
