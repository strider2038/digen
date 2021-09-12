package di_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strider2038/digen/pkg/di"
)

const testContainerSource = `
package di

import (
	"example.com/test/application/usecase"
	"example.com/test/domain"
	"example.com/test/di/config"
	httpadapter "example.com/test/infrastructure/api/http"
)

type Container struct {
	err error

	configuration    config.Configuration
	entityRepository domain.EntityRepository ` + "`di:\"required,set,close,public,external\"`" + `
	handler          *httpadapter.GetEntityHandler

	useCase *UseCaseContainer ` + "`di:\"container\"`" + `
}

type UseCaseContainer struct {
	*Container

	findEntity *usecase.FindEntity
}
`

const testFactorySource = `
package definitions

import (
	"example.com/test/application/usecase"
	"example.com/test/domain"
	httpadapter "example.com/test/infrastructure/api/http"
	"example.com/test/infrastructure/inmemory"
)

func CreateEntityRepository(c Container) domain.EntityRepository {
	return inmemory.NewEntityRepository()
}

func CreateUseCase(c Container) *usecase.FindEntity {
	return usecase.NewFindEntity(c.EntityRepository())
}

func CreateHandler(c Container) *httpadapter.GetEntityHandler {
	return httpadapter.NewGetEntityHandler(c.UseCase())
}
`

func TestParseContainerFromSource(t *testing.T) {
	container, err := di.ParseContainerFromSource(testContainerSource)

	require.NoError(t, err)
	require.NotNil(t, container)
	assert.Equal(t, "Container", container.Name)
	assert.Equal(t, "di", container.Package)
	assertExpectedContainerImports(t, container.Imports)
	assertExpectedContainerServices(t, container.Services)
	assertExpectedInternalContainers(t, container.Containers)
}

func TestParseFactoryFromSource(t *testing.T) {
	factory, err := di.ParseFactoryFromSource(testFactorySource)

	require.NoError(t, err)
	require.NotNil(t, factory)
	assert.NotNil(t, factory.Imports["usecase"])
	assert.NotNil(t, factory.Imports["domain"])
	assert.NotNil(t, factory.Imports["httpadapter"])
	assert.NotNil(t, factory.Imports["inmemory"])
	assert.Contains(t, factory.Services, "EntityRepository")
	assert.Contains(t, factory.Services, "UseCase")
	assert.Contains(t, factory.Services, "Handler")
}

func assertExpectedContainerImports(t *testing.T, imports map[string]*di.ImportDefinition) {
	if assert.NotNil(t, imports["usecase"]) {
		assert.Equal(t, "usecase", imports["usecase"].ID)
		assert.Equal(t, "", imports["usecase"].Name)
		assert.Equal(t, `"example.com/test/application/usecase"`, imports["usecase"].Path)
	}

	if assert.NotNil(t, imports["domain"]) {
		assert.Equal(t, "domain", imports["domain"].ID)
		assert.Equal(t, "", imports["domain"].Name)
		assert.Equal(t, `"example.com/test/domain"`, imports["domain"].Path)
	}

	if assert.NotNil(t, imports["config"]) {
		assert.Equal(t, "config", imports["config"].ID)
		assert.Equal(t, "", imports["config"].Name)
		assert.Equal(t, `"example.com/test/di/config"`, imports["config"].Path)
	}

	if assert.NotNil(t, imports["httpadapter"]) {
		assert.Equal(t, "httpadapter", imports["httpadapter"].ID)
		assert.Equal(t, "httpadapter", imports["httpadapter"].Name)
		assert.Equal(t, `"example.com/test/infrastructure/api/http"`, imports["httpadapter"].Path)
	}
}

func assertExpectedContainerServices(t *testing.T, services []*di.ServiceDefinition) {
	if assert.Len(t, services, 3) {
		assert.Equal(t, "configuration", services[0].Name)
		assert.Equal(t, "config", services[0].Type.Package)
		assert.Equal(t, "Configuration", services[0].Type.Name)
		assert.False(t, services[0].Type.IsPointer)
		assert.False(t, services[0].HasSetter)
		assert.False(t, services[0].HasCloser)
		assert.False(t, services[0].IsRequired)
		assert.False(t, services[0].IsPublic)
		assert.False(t, services[0].IsExternal)

		assert.Equal(t, "entityRepository", services[1].Name)
		assert.Equal(t, "domain", services[1].Type.Package)
		assert.Equal(t, "EntityRepository", services[1].Type.Name)
		assert.False(t, services[1].Type.IsPointer)
		assert.True(t, services[1].HasSetter)
		assert.True(t, services[1].HasCloser)
		assert.True(t, services[1].IsRequired)
		assert.True(t, services[1].IsPublic)
		assert.True(t, services[1].IsExternal)

		assert.Equal(t, "handler", services[2].Name)
		assert.Equal(t, "httpadapter", services[2].Type.Package)
		assert.Equal(t, "GetEntityHandler", services[2].Type.Name)
		assert.True(t, services[2].Type.IsPointer)
		assert.False(t, services[2].HasSetter)
		assert.False(t, services[2].HasCloser)
		assert.False(t, services[2].IsRequired)
		assert.False(t, services[2].IsPublic)
		assert.False(t, services[2].IsExternal)
	}
}

func assertExpectedInternalContainers(t *testing.T, containers []*di.ContainerDefinition) {
	if !assert.Len(t, containers, 1) {
		return
	}

	services := containers[0].Services
	if !assert.Len(t, services, 1) {
		return
	}

	assert.Equal(t, "findEntity", services[0].Name)
	assert.Equal(t, "usecase", services[0].Type.Package)
	assert.Equal(t, "FindEntity", services[0].Type.Name)
	assert.True(t, services[0].Type.IsPointer)
	assert.False(t, services[0].HasSetter)
	assert.False(t, services[0].HasCloser)
	assert.False(t, services[0].IsRequired)
	assert.False(t, services[0].IsPublic)
	assert.False(t, services[0].IsExternal)
}
