package di

import (
	"github.com/muonsoft/errors"
)

type FileGenerator struct {
	fileLocator FileLocator
	container   *RootContainerDefinition
	params      GenerationParameters
}

func NewFileGenerator(
	fileLocator FileLocator,
	container *RootContainerDefinition,
	params GenerationParameters,
) *FileGenerator {
	return &FileGenerator{
		fileLocator: fileLocator,
		container:   container,
		params:      params,
	}
}

func (g *FileGenerator) GenerateFiles() ([]*File, error) {
	files := make([]*File, 0)

	generators := [...]func() (*File, error){
		NewInternalContainerGenerator(g.fileLocator, g.container, g.params).Generate,
		NewLookupContainerGenerator(g.fileLocator, g.container, g.params).Generate,
		NewPublicContainerGenerator(g.fileLocator, g.container, g.params).Generate,
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
