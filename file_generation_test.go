package digen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strider2038/digen"
)

const basicContainerNewFile = `package testpkg

func NewContainer() *Container {
	c := &Container{}

	return c
}
`

const basicContainerGetFile = `package testpkg

import (
	"example.com/test/domain"
	"example.com/test/di/internal/definitions"
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
)

func (c *Container) SetServiceName(s *domain.Service) {
	c.serviceName = s
}
`

const basicContainerCloseFile = `package testpkg

func (c *Container) Close() {
	if c.connection != nil {
		c.connection.Close()
	}
}
`

const basicContainerDefinitionsContractsFile = `package definitions

import (
	"example.com/test/domain"
)

type Container interface {
	SetError(err error)

	ServiceName() *domain.Service
}
`

const basicContainerPublicFile = `package di

import (
	"sync"
	"example.com/test/di/internal"
	"example.com/test/domain"
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
	"example.com/test/di/internal"
	"example.com/test/domain"
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
	"example.com/test/di/internal"
	"example.com/test/domain"
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

const twoServicesContainerGetFile = `package testpkg

import (
	"example.com/test/domain"
	"example.com/test/di/internal/definitions"
)

func (c *Container) FirstService() *domain.Service {
	if c.firstService == nil {
		c.firstService = definitions.CreateFirstService(c)
	}
	return c.firstService
}

func (c *Container) SecondService() *domain.Service {
	if c.secondService == nil {
		c.secondService = definitions.CreateSecondService(c)
	}
	return c.secondService
}
`

const internalContainerNewFile = `package testpkg

func NewContainer() *Container {
	c := &Container{}
	c.internalContainerName = &InternalContainerType{Container: c}

	return c
}
`

const internalContainerGetFile = `package testpkg

import (
	"example.com/test/domain"
	"example.com/test/di/internal/definitions"
)

func (c *Container) TopService() *domain.Service {
	if c.topService == nil {
		c.topService = definitions.CreateTopService(c)
	}
	return c.topService
}

func (c *Container) InternalContainerName() definitions.InternalContainerType {
	return c.internalContainerName
}

func (c *InternalContainerType) FirstService() *domain.Service {
	if c.firstService == nil {
		c.firstService = definitions.CreateFirstService(c)
	}
	return c.firstService
}

func (c *InternalContainerType) SecondService() *domain.Service {
	if c.secondService == nil {
		c.secondService = definitions.CreateSecondService(c)
	}
	return c.secondService
}
`

const internalContainerSetFile = `package testpkg

import (
	"example.com/test/domain"
)

func (c *InternalContainerType) SetSecondService(s *domain.Service) {
	c.secondService = s
}
`

const internalContainerCloseFile = `package testpkg

func (c *Container) Close() {
	if c.internalContainerName.secondService != nil {
		c.internalContainerName.secondService.Close()
	}
}
`

const internalContainerDefinitionsContractsFile = `package definitions

import (
	"example.com/test/domain"
)

type Container interface {
	SetError(err error)

	TopService() *domain.Service

	InternalContainerName() InternalContainerType
}

type InternalContainerType interface {
	FirstService() *domain.Service
	SecondService() *domain.Service
}
`

const internalContainerPublicFile = `package di

import (
	"sync"
	"example.com/test/di/internal"
	"example.com/test/domain"
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

func (c *Container) FirstService() (*domain.Service, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	s := c.c.InternalContainerName().(*testpkg.InternalContainerType).FirstService()
	err := c.c.Error()
	if err != nil {
		return nil, err
	}

	return s, err
}

func (c *Container) SetSecondService(s *domain.Service) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.c.InternalContainerName().(*testpkg.InternalContainerType).SetSecondService(s)
}

func (c *Container) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.c.Close()
}
`

const factoryFile = `package definitions

import (
	"example.com/test/domain"
)

func CreateFirstService(c Container) *domain.Service {
	panic("not implemented")
}

func CreateSecondService(c Container) *domain.Service {
	panic("not implemented")
}
`

func TestGenerate(t *testing.T) {
	tests := []struct {
		name      string
		container *digen.RootContainerDefinition
		assert    func(t *testing.T, files []*digen.File)
	}{
		{
			name: "only basic getters",
			container: &digen.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*digen.ImportDefinition{
					"domain":      {Path: `"example.com/test/domain"`},
					"httpadapter": {Name: "httpadapter", Path: `"example.com/test/infrastructure/api/http"`},
				},
				Services: []*digen.ServiceDefinition{
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

				require.Len(t, files, 5)
				assert.Equal(t, digen.InternalPackage, files[0].Package)
				assert.Equal(t, "container_new.go", files[0].Name)
				assert.Equal(t, basicContainerNewFile, string(files[0].Content))
				assert.Equal(t, digen.InternalPackage, files[1].Package)
				assert.Equal(t, "container_get.go", files[1].Name)
				assert.Equal(t, basicContainerGetFile, string(files[1].Content))
				assert.Equal(t, digen.DefinitionsPackage, files[3].Package)
				assert.Equal(t, "contracts.go", files[3].Name)
				assert.Equal(t, basicContainerDefinitionsContractsFile, string(files[3].Content))
				assert.Equal(t, digen.PublicPackage, files[4].Package)
				assert.Equal(t, "container.go", files[4].Name)
				assert.Equal(t, basicContainerPublicFile, string(files[4].Content))
			},
		},
		{
			name: "definition with setter",
			container: &digen.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*digen.ImportDefinition{
					"domain":      {Path: `"example.com/test/domain"`},
					"httpadapter": {Name: "httpadapter", Path: `"example.com/test/infrastructure/api/http"`},
				},
				Services: []*digen.ServiceDefinition{
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

				require.GreaterOrEqual(t, len(files), 6)
				assert.Equal(t, digen.InternalPackage, files[2].Package)
				assert.Equal(t, "container_set.go", files[2].Name)
				assert.Equal(t, basicContainerSetFile, string(files[2].Content))
				assert.Equal(t, digen.PublicPackage, files[5].Package)
				assert.Equal(t, "container.go", files[5].Name)
				assert.Equal(t, publicContainerWithSetterFile, string(files[5].Content))
			},
		},
		{
			name: "definition with requirement",
			container: &digen.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*digen.ImportDefinition{
					"domain":      {Path: `"example.com/test/domain"`},
					"httpadapter": {Name: "httpadapter", Path: `"example.com/test/infrastructure/api/http"`},
				},
				Services: []*digen.ServiceDefinition{
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

				require.GreaterOrEqual(t, len(files), 6)
				assert.Equal(t, basicContainerSetFile, string(files[2].Content))
				assert.Equal(t, publicContainerWithRequirementFile, string(files[5].Content))
			},
		},
		{
			name: "definition external service",
			container: &digen.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*digen.ImportDefinition{
					"domain":      {Path: `"example.com/test/domain"`},
					"httpadapter": {Name: "httpadapter", Path: `"example.com/test/infrastructure/api/http"`},
				},
				Services: []*digen.ServiceDefinition{
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

				require.GreaterOrEqual(t, len(files), 3)
				assert.Equal(t, containerGetFileWithExternalService, string(files[1].Content))
				assert.Equal(t, basicContainerSetFile, string(files[2].Content))
			},
		},
		{
			name: "definition of static type",
			container: &digen.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*digen.ImportDefinition{
					"config": {Path: `"example.com/test/di/config"`},
				},
				Services: []*digen.ServiceDefinition{
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

				require.GreaterOrEqual(t, len(files), 2)
				assert.Equal(t, containerGetFileWithStaticType, string(files[1].Content))
			},
		},
		{
			name: "definition with closer",
			container: &digen.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*digen.ImportDefinition{
					"sql": {Path: `"example.com/test/sql"`},
				},
				Services: []*digen.ServiceDefinition{
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

				require.Greater(t, len(files), 3)
				assert.Equal(t, digen.InternalPackage, files[2].Package)
				assert.Equal(t, "container_close.go", files[2].Name)
				assert.Equal(t, basicContainerCloseFile, string(files[2].Content))
			},
		},
		{
			name: "two services from one package",
			container: &digen.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*digen.ImportDefinition{
					"domain": {Path: `"example.com/test/domain"`},
				},
				Services: []*digen.ServiceDefinition{
					{
						Name: "firstService",
						Type: digen.TypeDefinition{
							IsPointer: true,
							Package:   "domain",
							Name:      "Service",
						},
						IsPublic: true,
					},
					{
						Name: "secondService",
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

				require.Greater(t, len(files), 2)
				assert.Equal(t, twoServicesContainerGetFile, string(files[1].Content))
			},
		},
		{
			name: "internal container",
			container: &digen.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*digen.ImportDefinition{
					"domain": {Path: `"example.com/test/domain"`},
				},
				Services: []*digen.ServiceDefinition{
					{
						Name: "topService",
						Type: digen.TypeDefinition{
							IsPointer: true,
							Package:   "domain",
							Name:      "Service",
						},
					},
				},
				Containers: []*digen.ContainerDefinition{
					{
						Name: "internalContainerName",
						Type: digen.TypeDefinition{
							IsPointer: true,
							Package:   "testpkg",
							Name:      "InternalContainerType",
						},
						Services: []*digen.ServiceDefinition{
							{
								Name: "firstService",
								Type: digen.TypeDefinition{
									IsPointer: true,
									Package:   "domain",
									Name:      "Service",
								},
								IsPublic: true,
							},
							{
								Name: "secondService",
								Type: digen.TypeDefinition{
									IsPointer: true,
									Package:   "domain",
									Name:      "Service",
								},
								HasSetter: true,
								HasCloser: true,
							},
						},
					},
				},
			},
			assert: func(t *testing.T, files []*digen.File) {
				t.Helper()

				require.Greater(t, len(files), 5)
				assert.Equal(t, internalContainerNewFile, string(files[0].Content))
				assert.Equal(t, internalContainerGetFile, string(files[1].Content))
				assert.Equal(t, internalContainerSetFile, string(files[2].Content))
				assert.Equal(t, internalContainerCloseFile, string(files[3].Content))
				assert.Equal(t, internalContainerDefinitionsContractsFile, string(files[4].Content))
				assert.Equal(t, internalContainerPublicFile, string(files[5].Content))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			files, err := digen.GenerateFiles(test.container, digen.GenerationParameters{
				RootPackage: "example.com/test/di",
			})

			require.NoError(t, err)
			test.assert(t, files)
		})
	}
}

func TestGenerateFactory(t *testing.T) {
	container := &digen.RootContainerDefinition{
		Name:    "Container",
		Package: "testpkg",
		Imports: map[string]*digen.ImportDefinition{
			"domain":   {Path: `"example.com/test/domain"`},
			"external": {Path: `"example.com/test/external"`},
		},
		Services: []*digen.ServiceDefinition{
			{
				Name: "firstService",
				Type: digen.TypeDefinition{
					IsPointer: true,
					Package:   "domain",
					Name:      "Service",
				},
			},
			{
				Name: "secondService",
				Type: digen.TypeDefinition{
					IsPointer: true,
					Package:   "domain",
					Name:      "Service",
				},
			},
			{
				Name: "externalService",
				Type: digen.TypeDefinition{
					IsPointer: true,
					Package:   "external",
					Name:      "Service",
				},
				IsExternal: true,
			},
			{
				Name: "requiredService",
				Type: digen.TypeDefinition{
					IsPointer: true,
					Package:   "required",
					Name:      "Service",
				},
				IsRequired: true,
			},
		},
	}

	file, err := digen.GenerateFactory(container, digen.GenerationParameters{
		RootPackage: "example.com/test/di",
	})

	require.NoError(t, err)
	require.NotNil(t, file)
	assert.Equal(t, factoryFile, string(file.Content))
}
