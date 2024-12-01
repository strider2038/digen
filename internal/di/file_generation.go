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
	Package      string
	WrapPackage  string
	WrapFunction string
	Verb         string
}

func (w ErrorHandling) Defaults() ErrorHandling {
	if w.Package == "" {
		w.Package = "errors"
	}
	if w.WrapPackage == "" {
		w.WrapPackage = "fmt"
	}
	if w.WrapFunction == "" {
		w.WrapFunction = "Errorf"
	}
	if w.Verb == "" {
		w.Verb = "%w"
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
	path := params.ErrorHandling.WrapPackage
	funcName := params.ErrorHandling.WrapFunction
	verb := params.ErrorHandling.Verb

	return jen.Qual(path, funcName).Call(jen.Lit(message+": "+verb), errorIdentifier)
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
