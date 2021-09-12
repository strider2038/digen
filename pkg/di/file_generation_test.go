package di_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strider2038/digen/pkg/di"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name      string
		container *di.RootContainerDefinition
		assert    func(t *testing.T, files []*di.File)
	}{
		{
			name: "single container with getters only",
			container: &di.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*di.ImportDefinition{
					"domain":      {Path: `"example.com/test/domain"`},
					"httpadapter": {Name: "httpadapter", Path: `"example.com/test/infrastructure/api/http"`},
				},
				Services: []*di.ServiceDefinition{
					{
						Name: "serviceName",
						Type: di.TypeDefinition{
							IsPointer: true,
							Package:   "domain",
							Name:      "Service",
						},
						IsPublic: true,
					},
				},
			},
			assert: func(t *testing.T, files []*di.File) {
				t.Helper()

				require.Len(t, files, 3)
				assert.Equal(t, di.InternalPackage, files[0].Package)
				assert.Equal(t, "generated.go", files[0].Name)
				assert.Equal(t, singleContainerWithGettersOnlyGeneratedFile, string(files[0].Content))
				assert.Equal(t, di.DefinitionsPackage, files[1].Package)
				assert.Equal(t, "contracts.go", files[1].Name)
				assert.Equal(t, singleContainerWithGettersOnlyDefinitionContracts, string(files[1].Content))
				assert.Equal(t, di.PublicPackage, files[2].Package)
				assert.Equal(t, "container.go", files[2].Name)
				assert.Equal(t, singleContainerWithGettersOnlyPublicFile, string(files[2].Content))
			},
		},
		{
			name: "single container with service setter",
			container: &di.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*di.ImportDefinition{
					"domain":      {Path: `"example.com/test/domain"`},
					"httpadapter": {Name: "httpadapter", Path: `"example.com/test/infrastructure/api/http"`},
				},
				Services: []*di.ServiceDefinition{
					{
						Name: "serviceName",
						Type: di.TypeDefinition{
							IsPointer: true,
							Package:   "domain",
							Name:      "Service",
						},
						HasSetter: true,
					},
				},
			},
			assert: func(t *testing.T, files []*di.File) {
				t.Helper()

				require.GreaterOrEqual(t, len(files), 3)
				assert.Equal(t, di.InternalPackage, files[0].Package)
				assert.Equal(t, "generated.go", files[0].Name)
				assert.Equal(t, singleContainerWithServiceSetterGeneratedFile, string(files[0].Content))
				assert.Equal(t, di.PublicPackage, files[2].Package)
				assert.Equal(t, "container.go", files[2].Name)
				assert.Equal(t, singleContainerWithServiceSetterPublicContainer, string(files[2].Content))
			},
		},
		{
			name: "single container with required service",
			container: &di.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*di.ImportDefinition{
					"domain":      {Path: `"example.com/test/domain"`},
					"httpadapter": {Name: "httpadapter", Path: `"example.com/test/infrastructure/api/http"`},
				},
				Services: []*di.ServiceDefinition{
					{
						Name: "serviceName",
						Type: di.TypeDefinition{
							IsPointer: true,
							Package:   "domain",
							Name:      "Service",
						},
						IsRequired: true,
					},
				},
			},
			assert: func(t *testing.T, files []*di.File) {
				t.Helper()

				require.GreaterOrEqual(t, len(files), 3)
				assert.Equal(t, singleContainerWithRequiredServiceGeneratedFile, string(files[0].Content))
				assert.Equal(t, publicContainerWithRequirementFile, string(files[2].Content))
			},
		},
		{
			name: "single container with external service",
			container: &di.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*di.ImportDefinition{
					"domain":      {Path: `"example.com/test/domain"`},
					"httpadapter": {Name: "httpadapter", Path: `"example.com/test/infrastructure/api/http"`},
				},
				Services: []*di.ServiceDefinition{
					{
						Name: "serviceName",
						Type: di.TypeDefinition{
							IsPointer: true,
							Package:   "domain",
							Name:      "Service",
						},
						IsExternal: true,
					},
				},
			},
			assert: func(t *testing.T, files []*di.File) {
				t.Helper()

				require.GreaterOrEqual(t, len(files), 1)
				assert.Equal(t, singleContainerWithExternalServiceGeneratedFile, string(files[0].Content))
			},
		},
		{
			name: "single container with static type",
			container: &di.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*di.ImportDefinition{
					"config": {Path: `"example.com/test/di/config"`},
				},
				Services: []*di.ServiceDefinition{
					{
						Name: "configuration",
						Type: di.TypeDefinition{
							Package: "config",
							Name:    "Configuration",
						},
						IsRequired: true,
					},
				},
			},
			assert: func(t *testing.T, files []*di.File) {
				t.Helper()

				require.GreaterOrEqual(t, len(files), 1)
				assert.Equal(t, singleContainerWithStaticTypeGeneratedFile, string(files[0].Content))
			},
		},
		{
			name: "single container with closer",
			container: &di.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*di.ImportDefinition{
					"sql": {Path: `"example.com/test/sql"`},
				},
				Services: []*di.ServiceDefinition{
					{
						Name: "connection",
						Type: di.TypeDefinition{
							Package: "sql",
							Name:    "Connection",
						},
						HasCloser: true,
					},
				},
			},
			assert: func(t *testing.T, files []*di.File) {
				t.Helper()

				require.Greater(t, len(files), 1)
				assert.Equal(t, di.InternalPackage, files[0].Package)
				assert.Equal(t, singleContainerWithCloserGeneratedFile, string(files[0].Content))
			},
		},
		{
			name: "two services from one package",
			container: &di.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*di.ImportDefinition{
					"domain": {Path: `"example.com/test/domain"`},
				},
				Services: []*di.ServiceDefinition{
					{
						Name: "firstService",
						Type: di.TypeDefinition{
							IsPointer: true,
							Package:   "domain",
							Name:      "Service",
						},
						IsPublic: true,
					},
					{
						Name: "secondService",
						Type: di.TypeDefinition{
							IsPointer: true,
							Package:   "domain",
							Name:      "Service",
						},
						IsPublic: true,
					},
				},
			},
			assert: func(t *testing.T, files []*di.File) {
				t.Helper()

				require.Greater(t, len(files), 1)
				assert.Equal(t, twoServicesFromOnePackageGeneratedFile, string(files[0].Content))
			},
		},
		{
			name: "separate container",
			container: &di.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*di.ImportDefinition{
					"domain": {Path: `"example.com/test/domain"`},
				},
				Services: []*di.ServiceDefinition{
					{
						Name: "topService",
						Type: di.TypeDefinition{
							IsPointer: true,
							Package:   "domain",
							Name:      "Service",
						},
					},
				},
				Containers: []*di.ContainerDefinition{
					{
						Name: "internalContainerName",
						Type: di.TypeDefinition{
							IsPointer: true,
							Package:   "testpkg",
							Name:      "InternalContainerType",
						},
						Services: []*di.ServiceDefinition{
							{
								Name: "firstService",
								Type: di.TypeDefinition{
									IsPointer: true,
									Package:   "domain",
									Name:      "Service",
								},
								IsPublic: true,
							},
							{
								Name: "secondService",
								Type: di.TypeDefinition{
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
			assert: func(t *testing.T, files []*di.File) {
				t.Helper()

				require.GreaterOrEqual(t, len(files), 3)
				assert.Equal(t, separateContainerGeneratedFile, string(files[0].Content))
				assert.Equal(t, separateContainerDefinitionsContractsFile, string(files[1].Content))
				assert.Equal(t, separateContainerPublicFile, string(files[2].Content))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			files, err := di.GenerateFiles(test.container, di.GenerationParameters{
				RootPackage: "example.com/test/di",
			})

			require.NoError(t, err)
			test.assert(t, files)
		})
	}
}

const singleContainerWithGettersOnlyGeneratedFile = `package testpkg

import (
	"example.com/test/domain"
	"example.com/test/di/internal/definitions"
)

func NewContainer() *Container {
	c := &Container{}

	return c
}

func (c *Container) ServiceName() *domain.Service {
	if c.serviceName == nil {
		c.serviceName = definitions.CreateServiceName(c)
	}
	return c.serviceName
}

func (c *Container) Close() {}
`

const singleContainerWithGettersOnlyDefinitionContracts = `package definitions

import (
	"example.com/test/domain"
)

type Container interface {
	SetError(err error)

	ServiceName() *domain.Service
}
`

const singleContainerWithGettersOnlyPublicFile = `package di

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

const singleContainerWithServiceSetterGeneratedFile = `package testpkg

import (
	"example.com/test/domain"
	"example.com/test/di/internal/definitions"
)

func NewContainer() *Container {
	c := &Container{}

	return c
}

func (c *Container) ServiceName() *domain.Service {
	if c.serviceName == nil {
		c.serviceName = definitions.CreateServiceName(c)
	}
	return c.serviceName
}

func (c *Container) SetServiceName(s *domain.Service) {
	c.serviceName = s
}

func (c *Container) Close() {}
`

const singleContainerWithServiceSetterPublicContainer = `package di

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

const singleContainerWithRequiredServiceGeneratedFile = `package testpkg

import (
	"example.com/test/domain"
)

func NewContainer() *Container {
	c := &Container{}

	return c
}

func (c *Container) ServiceName() *domain.Service {
	return c.serviceName
}

func (c *Container) SetServiceName(s *domain.Service) {
	c.serviceName = s
}

func (c *Container) Close() {}
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

const singleContainerWithExternalServiceGeneratedFile = `package testpkg

import (
	"example.com/test/domain"
)

func NewContainer() *Container {
	c := &Container{}

	return c
}

func (c *Container) ServiceName() *domain.Service {
	if c.serviceName == nil {
		panic("missing ServiceName")
	}
	return c.serviceName
}

func (c *Container) SetServiceName(s *domain.Service) {
	c.serviceName = s
}

func (c *Container) Close() {}
`

const singleContainerWithStaticTypeGeneratedFile = `package testpkg

import (
	"example.com/test/di/config"
)

func NewContainer() *Container {
	c := &Container{}

	return c
}

func (c *Container) Configuration() config.Configuration {
	return c.configuration
}

func (c *Container) SetConfiguration(s config.Configuration) {
	c.configuration = s
}

func (c *Container) Close() {}
`

const singleContainerWithCloserGeneratedFile = `package testpkg

import (
	"example.com/test/sql"
	"example.com/test/di/internal/definitions"
)

func NewContainer() *Container {
	c := &Container{}

	return c
}

func (c *Container) Connection() sql.Connection {
	if c.connection == nil {
		c.connection = definitions.CreateConnection(c)
	}
	return c.connection
}

func (c *Container) Close() {
	if c.connection != nil {
		c.connection.Close()
	}
}
`

const twoServicesFromOnePackageGeneratedFile = `package testpkg

import (
	"example.com/test/domain"
	"example.com/test/di/internal/definitions"
)

func NewContainer() *Container {
	c := &Container{}

	return c
}

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

func (c *Container) Close() {}
`

const separateContainerGeneratedFile = `package testpkg

import (
	"example.com/test/domain"
	"example.com/test/di/internal/definitions"
)

func NewContainer() *Container {
	c := &Container{}
	c.internalContainerName = &InternalContainerType{Container: c}

	return c
}

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

func (c *InternalContainerType) SetSecondService(s *domain.Service) {
	c.secondService = s
}

func (c *Container) Close() {
	if c.internalContainerName.secondService != nil {
		c.internalContainerName.secondService.Close()
	}
}
`

const separateContainerDefinitionsContractsFile = `package definitions

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

const separateContainerPublicFile = `package di

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
