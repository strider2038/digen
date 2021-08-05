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

const basicPublicContainerFile = `package di

import (
	"sync"

	"example.com/test/domain"
	httpadapter "example.com/test/infrastructure/api/http"
)

type Container struct {
	mu *sync.Mutex
	c  *internal.Container
}

type Injector func(c *Container) error

func NewContainer(injectors ...Injector) (*Container, error) {
	c := &Container{
		mu: &sync.Mutex{},
		c:  &internal.Container{},
	}

	for _, inject := range injectors {
		err := inject(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Container) ServiceName() (*domain.Service, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	s := c.c.ServiceName()
	err := c.c.Error()
	if err != nil {
		return nil, err
	}

	return s, err
}

func (c *Container) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.c.Close()
}
`

const publicContainerWithSetterFile = `package di

import (
	"sync"

	"example.com/test/domain"
	httpadapter "example.com/test/infrastructure/api/http"
)

type Container struct {
	mu *sync.Mutex
	c  *internal.Container
}

type Injector func(c *Container) error

func NewContainer(injectors ...Injector) (*Container, error) {
	c := &Container{
		mu: &sync.Mutex{},
		c:  &internal.Container{},
	}

	for _, inject := range injectors {
		err := inject(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Container) SetServiceName(s *domain.Service) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.c.SetServiceName(s)
}

func (c *Container) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.c.Close()
}
`

const publicContainerWithRequirementFile = `package di

import (
	"sync"

	"example.com/test/domain"
	httpadapter "example.com/test/infrastructure/api/http"
)

type Container struct {
	mu *sync.Mutex
	c  *internal.Container
}

type Injector func(c *Container) error

func NewContainer(
	serviceName *domain.Service,
	injectors ...Injector,
) (*Container, error) {
	c := &Container{
		mu: &sync.Mutex{},
		c:  &internal.Container{},
	}

	c.c.SetServiceName(serviceName)

	for _, inject := range injectors {
		err := inject(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Container) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.c.Close()
}
`

const containerGetFileWithExternalService = `package testpkg

import (
	"example.com/test/domain"
	httpadapter "example.com/test/infrastructure/api/http"
)

func (c *Container) ServiceName() *domain.Service {
	if c.serviceName == nil {
		panic("missing ServiceName")
	}
	return c.serviceName
}
`

const containerGetFileWithStaticType = `package testpkg

import (
	"example.com/test/di/config"
)

func (c *Container) Configuration() config.Configuration {
	return c.configuration
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
				Name:    "Container",
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
						IsPublic: true,
					},
				},
			},
			assert: func(t *testing.T, files []*digen.File) {
				t.Helper()

				require.Len(t, files, 4)
				assert.Equal(t, digen.InternalPackage, files[0].Package)
				assert.Equal(t, "container_get.go", files[0].Name)
				assert.Equal(t, basicContainerGetFile, string(files[0].Content))
				assert.Equal(t, digen.DefinitionsPackage, files[2].Package)
				assert.Equal(t, "definitions/container.go", files[2].Name)
				assert.Equal(t, basicDefinitionsContainerFile, string(files[2].Content))
				assert.Equal(t, digen.PublicPackage, files[3].Package)
				assert.Equal(t, "container.go", files[3].Name)
				assert.Equal(t, basicPublicContainerFile, string(files[3].Content))
			},
		},
		{
			name: "definition with setter",
			container: &digen.ContainerDefinition{
				Name:    "Container",
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

				require.GreaterOrEqual(t, len(files), 5)
				assert.Equal(t, digen.InternalPackage, files[1].Package)
				assert.Equal(t, "container_set.go", files[1].Name)
				assert.Equal(t, basicContainerSetFile, string(files[1].Content))
				assert.Equal(t, digen.PublicPackage, files[4].Package)
				assert.Equal(t, "container.go", files[4].Name)
				assert.Equal(t, publicContainerWithSetterFile, string(files[4].Content))
			},
		},
		{
			name: "definition with requirement",
			container: &digen.ContainerDefinition{
				Name:    "Container",
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
						IsRequired: true,
					},
				},
			},
			assert: func(t *testing.T, files []*digen.File) {
				t.Helper()

				require.GreaterOrEqual(t, len(files), 5)
				assert.Equal(t, basicContainerSetFile, string(files[1].Content))
				assert.Equal(t, publicContainerWithRequirementFile, string(files[4].Content))
			},
		},
		{
			name: "definition external service",
			container: &digen.ContainerDefinition{
				Name:    "Container",
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
						IsExternal: true,
					},
				},
			},
			assert: func(t *testing.T, files []*digen.File) {
				t.Helper()

				require.GreaterOrEqual(t, len(files), 2)
				assert.Equal(t, containerGetFileWithExternalService, string(files[0].Content))
				assert.Equal(t, basicContainerSetFile, string(files[1].Content))
			},
		},
		{
			name: "definition of static type",
			container: &digen.ContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: []digen.ImportDefinition{
					{Path: `"example.com/test/di/config"`},
				},
				Services: []digen.ServiceDefinition{
					{
						Name: "configuration",
						Type: digen.TypeDefinition{
							Package: "config",
							Name:    "Configuration",
						},
						IsRequired: true,
					},
				},
			},
			assert: func(t *testing.T, files []*digen.File) {
				t.Helper()

				require.GreaterOrEqual(t, len(files), 1)
				assert.Equal(t, containerGetFileWithStaticType, string(files[0].Content))
			},
		},
		{
			name: "definition with closer",
			container: &digen.ContainerDefinition{
				Name:    "Container",
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
