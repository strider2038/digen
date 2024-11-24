package di_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strider2038/digen/internal/di"
)

func TestFactoriesGenerator_Generate(t *testing.T) {
	tests := []struct {
		name      string
		container *di.RootContainerDefinition
		assert    func(t *testing.T, files []*di.File)
	}{
		{
			name: "container with services",
			container: &di.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*di.ImportDefinition{
					"domain": {Path: "example.com/test/domain"},
				},
				Services: []*di.ServiceDefinition{
					{
						Name: "serviceName",
						Type: di.TypeDefinition{
							IsPointer: true,
							Package:   "domain",
							Name:      "Service",
						},
					},
				},
			},
			assert: func(t *testing.T, files []*di.File) {
				t.Helper()

				require.Len(t, files, 1)
				assert.Equal(t, di.FactoriesPackage, files[0].Package)
				assert.Equal(t, "container.go", files[0].Name)
				assert.Equal(t, containerWithServices, string(files[0].Content))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			generator := di.NewFactoriesGenerator(test.container, "", di.GenerationParameters{
				RootPackage: "example.com/test/di",
			})
			files, err := generator.Generate()

			require.NoError(t, err)
			test.assert(t, files)
		})
	}
}

var (
	//go:embed testdata/generation/factories_container_with_services.txt
	containerWithServices string
)
