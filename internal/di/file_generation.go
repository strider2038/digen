package di

import (
	"github.com/muonsoft/errors"
)

func generateDefinitionsContainerFile() *File {
	return &File{
		Package: DefinitionsPackage,
		Name:    "container.go",
		Content: []byte(definitionsContainerFileSkeleton),
	}
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
		NewLookupContainerGenerator(g.container, g.params).Generate,
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
