package digen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strider2038/digen"
)

const basicContainerGetFile = `package testpkg

import (
	"example.com/test/domain"
	httpadapter "example.com/test/infrastructure/api/http"
)

func (c *Container) ServiceName() *domain.Service {
	if c.serviceName == nil {
		c.serviceName = definitions.CreateServiceName(c)
	}
	return c.serviceName
}
`

const basicContainerSetFile = `package testpkg

import (
	"example.com/test/domain"
	httpadapter "example.com/test/infrastructure/api/http"
)

func (c *Container) SetServiceName(s *domain.Service) {
	c.serviceName = s
}
`

const basicContainerCloseFile = `package testpkg

import (
	"example.com/test/sql"
)

func (c *Container) Close() {
	if c.connection != nil {
		c.connection.Close()
	}
}
`

const basicDefinitionsContainerFile = `package definitions

import (
	"example.com/test/domain"
	httpadapter "example.com/test/infrastructure/api/http"
)

type Container interface {
	ServiceName() *domain.Service
}
`

func TestGenerate(t *testing.T) {
	tests := []struct {
		name      string
		container *digen.ContainerDefinition
		assert    func(t *testing.T, files []*digen.File)
	}{
		{
			name: "only basic getters",
			container: &digen.ContainerDefinition{
				Package: "testpkg",
				Imports: []digen.ImportDefinition{
					{Path: `"example.com/test/domain"`},
					{Name: "httpadapter", Path: `"example.com/test/infrastructure/api/http"`},
				},
				Services: []digen.ServiceDefinition{
					{
						Name: "serviceName",
						Type: digen.TypeDefinition{
							IsPointer: true,
							Package:   "domain",
							Name:      "Service",
						},
					},
				},
			},
			assert: func(t *testing.T, files []*digen.File) {
				t.Helper()

				require.Len(t, files, 3)
				assert.Equal(t, digen.InternalPackage, files[0].Package)
				assert.Equal(t, "container_get.go", files[0].Name)
				assert.Equal(t, basicContainerGetFile, string(files[0].Content))
				assert.Equal(t, digen.DefinitionsPackage, files[2].Package)
				assert.Equal(t, "definitions/container.go", files[2].Name)
				assert.Equal(t, basicDefinitionsContainerFile, string(files[2].Content))
			},
		},
		{
			name: "definition with setter",
			container: &digen.ContainerDefinition{
				Package: "testpkg",
				Imports: []digen.ImportDefinition{
					{Path: `"example.com/test/domain"`},
					{Name: "httpadapter", Path: `"example.com/test/infrastructure/api/http"`},
				},
				Services: []digen.ServiceDefinition{
					{
						Name: "serviceName",
						Type: digen.TypeDefinition{
							IsPointer: true,
							Package:   "domain",
							Name:      "Service",
						},
						HasSetter: true,
					},
				},
			},
			assert: func(t *testing.T, files []*digen.File) {
				t.Helper()

				require.Greater(t, len(files), 2)
				assert.Equal(t, digen.InternalPackage, files[1].Package)
				assert.Equal(t, "container_set.go", files[1].Name)
				assert.Equal(t, basicContainerSetFile, string(files[1].Content))
			},
		},
		{
			name: "definition with closer",
			container: &digen.ContainerDefinition{
				Package: "testpkg",
				Imports: []digen.ImportDefinition{
					{Path: `"example.com/test/sql"`},
				},
				Services: []digen.ServiceDefinition{
					{
						Name: "connection",
						Type: digen.TypeDefinition{
							Package: "sql",
							Name:    "Connection",
						},
						HasCloser: true,
					},
				},
			},
			assert: func(t *testing.T, files []*digen.File) {
				t.Helper()

				require.Greater(t, len(files), 2)
				assert.Equal(t, digen.InternalPackage, files[1].Package)
				assert.Equal(t, "container_close.go", files[1].Name)
				assert.Equal(t, basicContainerCloseFile, string(files[1].Content))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			files, err := digen.Generate(test.container, digen.GenerationParameters{})

			require.NoError(t, err)
			test.assert(t, files)
		})
	}
}
