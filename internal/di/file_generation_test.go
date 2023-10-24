package di_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strider2038/digen/internal/di"
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
				assert.Equal(t, "container.go", files[0].Name)
				assert.Equal(t, singleContainerWithGettersOnlyInternalContainer, string(files[0].Content))
				assert.Equal(t, di.LookupPackage, files[1].Package)
				assert.Equal(t, "container.go", files[1].Name)
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
				assert.Equal(t, "container.go", files[0].Name)
				assert.Equal(t, singleContainerWithServiceSetterInternalContainer, string(files[0].Content))
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
				assert.Equal(t, singleContainerWithRequiredServiceInternalContainer, string(files[0].Content))
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
				assert.Equal(t, singleContainerWithExternalServiceInternalContainer, string(files[0].Content))
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
				assert.Equal(t, singleContainerWithStaticTypeInternalContainer, string(files[0].Content))
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
				assert.Equal(t, singleContainerWithCloserInternalContainer, string(files[0].Content))
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
				assert.Equal(t, twoServicesFromOnePackageInternalContainer, string(files[0].Content))
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
							// IsPointer: true,
							Package: "testpkg",
							Name:    "InternalContainerType",
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
				assert.Equal(t, separateContainerInternalContainer, string(files[0].Content))
				assert.Equal(t, separateLookupContainerFile, string(files[1].Content))
				assert.Equal(t, separateContainerPublicFile, string(files[2].Content))
			},
		},
		{
			name: "override service public name",
			container: &di.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*di.ImportDefinition{
					"http": {Path: `"net/http"`},
				},
				Services: []*di.ServiceDefinition{
					{
						Name:       "Router",
						PublicName: "APIRouter",
						IsPublic:   true,
						Type: di.TypeDefinition{
							Package: "http",
							Name:    "Handler",
						},
					},
				},
			},
			assert: func(t *testing.T, files []*di.File) {
				t.Helper()

				fmt.Println(string(files[2].Content))
				require.GreaterOrEqual(t, len(files), 3)
				assert.Contains(
					t, string(files[2].Content),
					`func (c *Container) APIRouter(ctx context.Context) (http.Handler, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	s := c.c.Router(ctx)
	err := c.c.Error()
	if err != nil {
		return nil, err
	}

	return s, err
}`,
				)
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

const singleContainerWithGettersOnlyInternalContainer = `package internal

import (
	"context"
	"example.com/test/di/internal/factories"
	"example.com/test/domain"
)

type Container struct {
	err error

	serviceName *domain.Service
}

func NewContainer() *Container {
	c := &Container{}

	return c
}

// Error returns the first initialization error, which can be set via SetError in a service definition.
func (c *Container) Error() error {
	return c.err
}

// SetError sets the first error into container. The error is used in the public container to return an initialization error.
func (c *Container) SetError(err error) {
	if err != nil && c.err == nil {
		c.err = err
	}
}

func (c *Container) ServiceName(ctx context.Context) *domain.Service {
	if c.serviceName == nil && c.err == nil {
		c.serviceName = factories.CreateServiceName(ctx, c)
	}
	return c.serviceName
}

func (c *Container) Close() {}
`

const singleContainerWithGettersOnlyDefinitionContracts = `package lookup

import (
	"context"
	"example.com/test/domain"
)

type Container interface {
	// SetError sets the first error into container. The error is used in the public container to return an initialization error.
	SetError(err error)

	ServiceName(ctx context.Context) *domain.Service
}
`

const singleContainerWithGettersOnlyPublicFile = `package di

import (
	"context"
	"example.com/test/di/internal"
	"example.com/test/domain"
	"sync"
)

type Container struct {
	mu *sync.Mutex
	c  *internal.Container
}

type Injector func(c *Container) error

func NewContainer(injectors ...Injector) (*Container, error) {
	c := &Container{
		mu: &sync.Mutex{},
		c:  internal.NewContainer(),
	}

	for _, inject := range injectors {
		err := inject(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Container) ServiceName(ctx context.Context) (*domain.Service, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	s := c.c.ServiceName(ctx)
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

const singleContainerWithServiceSetterInternalContainer = `package internal

import (
	"context"
	"example.com/test/di/internal/factories"
	"example.com/test/domain"
)

type Container struct {
	err error

	serviceName *domain.Service
}

func NewContainer() *Container {
	c := &Container{}

	return c
}

// Error returns the first initialization error, which can be set via SetError in a service definition.
func (c *Container) Error() error {
	return c.err
}

// SetError sets the first error into container. The error is used in the public container to return an initialization error.
func (c *Container) SetError(err error) {
	if err != nil && c.err == nil {
		c.err = err
	}
}

func (c *Container) ServiceName(ctx context.Context) *domain.Service {
	if c.serviceName == nil && c.err == nil {
		c.serviceName = factories.CreateServiceName(ctx, c)
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
	"example.com/test/di/internal"
	"example.com/test/domain"
	"sync"
)

type Container struct {
	mu *sync.Mutex
	c  *internal.Container
}

type Injector func(c *Container) error

func NewContainer(injectors ...Injector) (*Container, error) {
	c := &Container{
		mu: &sync.Mutex{},
		c:  internal.NewContainer(),
	}

	for _, inject := range injectors {
		err := inject(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func SetServiceName(s *domain.Service) Injector {
	return func(c *Container) error {
		c.c.SetServiceName(s)

		return nil
	}
}

func (c *Container) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.c.Close()
}
`

const singleContainerWithRequiredServiceInternalContainer = `package internal

import (
	"context"
	"example.com/test/domain"
)

type Container struct {
	err error

	serviceName *domain.Service
}

func NewContainer() *Container {
	c := &Container{}

	return c
}

// Error returns the first initialization error, which can be set via SetError in a service definition.
func (c *Container) Error() error {
	return c.err
}

// SetError sets the first error into container. The error is used in the public container to return an initialization error.
func (c *Container) SetError(err error) {
	if err != nil && c.err == nil {
		c.err = err
	}
}

func (c *Container) ServiceName(ctx context.Context) *domain.Service {
	return c.serviceName
}

func (c *Container) SetServiceName(s *domain.Service) {
	c.serviceName = s
}

func (c *Container) Close() {}
`

const publicContainerWithRequirementFile = `package di

import (
	"example.com/test/di/internal"
	"example.com/test/domain"
	"sync"
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
		c:  internal.NewContainer(),
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

const singleContainerWithExternalServiceInternalContainer = `package internal

import (
	"context"
	"example.com/test/domain"
)

type Container struct {
	err error

	serviceName *domain.Service
}

func NewContainer() *Container {
	c := &Container{}

	return c
}

// Error returns the first initialization error, which can be set via SetError in a service definition.
func (c *Container) Error() error {
	return c.err
}

// SetError sets the first error into container. The error is used in the public container to return an initialization error.
func (c *Container) SetError(err error) {
	if err != nil && c.err == nil {
		c.err = err
	}
}

func (c *Container) ServiceName(ctx context.Context) *domain.Service {
	if c.serviceName == nil && c.err == nil {
		panic("missing ServiceName")
	}
	return c.serviceName
}

func (c *Container) SetServiceName(s *domain.Service) {
	c.serviceName = s
}

func (c *Container) Close() {}
`

const singleContainerWithStaticTypeInternalContainer = `package internal

import (
	"context"
	"example.com/test/di/config"
)

type Container struct {
	err error

	configuration config.Configuration
}

func NewContainer() *Container {
	c := &Container{}

	return c
}

// Error returns the first initialization error, which can be set via SetError in a service definition.
func (c *Container) Error() error {
	return c.err
}

// SetError sets the first error into container. The error is used in the public container to return an initialization error.
func (c *Container) SetError(err error) {
	if err != nil && c.err == nil {
		c.err = err
	}
}

func (c *Container) Configuration(ctx context.Context) config.Configuration {
	return c.configuration
}

func (c *Container) SetConfiguration(s config.Configuration) {
	c.configuration = s
}

func (c *Container) Close() {}
`

const singleContainerWithCloserInternalContainer = `package internal

import (
	"context"
	"example.com/test/di/internal/factories"
	"example.com/test/sql"
)

type Container struct {
	err error

	connection sql.Connection
}

func NewContainer() *Container {
	c := &Container{}

	return c
}

// Error returns the first initialization error, which can be set via SetError in a service definition.
func (c *Container) Error() error {
	return c.err
}

// SetError sets the first error into container. The error is used in the public container to return an initialization error.
func (c *Container) SetError(err error) {
	if err != nil && c.err == nil {
		c.err = err
	}
}

func (c *Container) Connection(ctx context.Context) sql.Connection {
	if c.connection == nil && c.err == nil {
		c.connection = factories.CreateConnection(ctx, c)
	}
	return c.connection
}

func (c *Container) Close() {
	if c.connection != nil {
		c.connection.Close()
	}
}
`

const twoServicesFromOnePackageInternalContainer = `package internal

import (
	"context"
	"example.com/test/di/internal/factories"
	"example.com/test/domain"
)

type Container struct {
	err error

	firstService  *domain.Service
	secondService *domain.Service
}

func NewContainer() *Container {
	c := &Container{}

	return c
}

// Error returns the first initialization error, which can be set via SetError in a service definition.
func (c *Container) Error() error {
	return c.err
}

// SetError sets the first error into container. The error is used in the public container to return an initialization error.
func (c *Container) SetError(err error) {
	if err != nil && c.err == nil {
		c.err = err
	}
}

func (c *Container) FirstService(ctx context.Context) *domain.Service {
	if c.firstService == nil && c.err == nil {
		c.firstService = factories.CreateFirstService(ctx, c)
	}
	return c.firstService
}

func (c *Container) SecondService(ctx context.Context) *domain.Service {
	if c.secondService == nil && c.err == nil {
		c.secondService = factories.CreateSecondService(ctx, c)
	}
	return c.secondService
}

func (c *Container) Close() {}
`

const separateContainerInternalContainer = `package internal

import (
	"context"
	"example.com/test/di/internal/factories"
	"example.com/test/di/internal/lookup"
	"example.com/test/domain"
)

type Container struct {
	err error

	topService *domain.Service

	internalContainerName *InternalContainerType
}

func NewContainer() *Container {
	c := &Container{}
	c.internalContainerName = &InternalContainerType{Container: c}

	return c
}

// Error returns the first initialization error, which can be set via SetError in a service definition.
func (c *Container) Error() error {
	return c.err
}

// SetError sets the first error into container. The error is used in the public container to return an initialization error.
func (c *Container) SetError(err error) {
	if err != nil && c.err == nil {
		c.err = err
	}
}

type InternalContainerType struct {
	*Container

	firstService  *domain.Service
	secondService *domain.Service
}

func (c *Container) TopService(ctx context.Context) *domain.Service {
	if c.topService == nil && c.err == nil {
		c.topService = factories.CreateTopService(ctx, c)
	}
	return c.topService
}

func (c *Container) InternalContainerName() lookup.InternalContainerType {
	return c.internalContainerName
}

func (c *InternalContainerType) FirstService(ctx context.Context) *domain.Service {
	if c.firstService == nil && c.err == nil {
		c.firstService = factories.CreateFirstService(ctx, c)
	}
	return c.firstService
}

func (c *InternalContainerType) SecondService(ctx context.Context) *domain.Service {
	if c.secondService == nil && c.err == nil {
		c.secondService = factories.CreateSecondService(ctx, c)
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

const separateLookupContainerFile = `package lookup

import (
	"context"
	"example.com/test/domain"
)

type Container interface {
	// SetError sets the first error into container. The error is used in the public container to return an initialization error.
	SetError(err error)

	TopService(ctx context.Context) *domain.Service

	InternalContainerName() InternalContainerType
}

type InternalContainerType interface {
	FirstService(ctx context.Context) *domain.Service
	SecondService(ctx context.Context) *domain.Service
}
`

const separateContainerPublicFile = `package di

import (
	"context"
	"example.com/test/di/internal"
	"example.com/test/domain"
	"sync"
)

type Container struct {
	mu *sync.Mutex
	c  *internal.Container
}

type Injector func(c *Container) error

func NewContainer(injectors ...Injector) (*Container, error) {
	c := &Container{
		mu: &sync.Mutex{},
		c:  internal.NewContainer(),
	}

	for _, inject := range injectors {
		err := inject(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Container) FirstService(ctx context.Context) (*domain.Service, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	s := c.c.InternalContainerName().(*testpkg.InternalContainerType).FirstService(ctx)
	err := c.c.Error()
	if err != nil {
		return nil, err
	}

	return s, err
}

func SetSecondService(s *domain.Service) Injector {
	return func(c *Container) error {
		c.c.InternalContainerName().(*testpkg.InternalContainerType).SetSecondService(s)

		return nil
	}
}

func (c *Container) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.c.Close()
}
`
