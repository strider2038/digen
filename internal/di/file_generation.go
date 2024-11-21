package di

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/muonsoft/errors"
)

func GenerateDefinitionsContainerFile() (*File, error) {
	var buffer bytes.Buffer

	err := definitionsContainerFileTemplate.Execute(&buffer, nil)
	if err != nil {
		return nil, errors.Errorf("generate internal container: %w")
	}

	return &File{
		Package: DefinitionsPackage,
		Name:    "container.go",
		Content: buffer.Bytes(),
	}, nil
}

type GenerationParameters struct {
	RootPackage string
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
