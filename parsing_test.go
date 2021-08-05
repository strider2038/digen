package digen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strider2038/digen"
)

const testSource = `
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
	useCase          *usecase.FindEntity
	handler          *httpadapter.GetEntityHandler
}
`

func TestParseSource(t *testing.T) {
	container, err := digen.ParseSource(testSource)

	require.NoError(t, err)
	require.NotNil(t, container)
	assert.Equal(t, "Container", container.Name)
	assert.Equal(t, "di", container.Package)
	assertExpectedImports(t, container.Imports)
	assertExpectedServices(t, container)
}

func assertExpectedImports(t *testing.T, imports []digen.ImportDefinition) {
	if assert.Len(t, imports, 4) {
		assert.Equal(t, "usecase", imports[0].ID)
		assert.Equal(t, "", imports[0].Name)
		assert.Equal(t, `"example.com/test/application/usecase"`, imports[0].Path)

		assert.Equal(t, "domain", imports[1].ID)
		assert.Equal(t, "", imports[1].Name)
		assert.Equal(t, `"example.com/test/domain"`, imports[1].Path)

		assert.Equal(t, "config", imports[2].ID)
		assert.Equal(t, "", imports[2].Name)
		assert.Equal(t, `"example.com/test/di/config"`, imports[2].Path)

		assert.Equal(t, "httpadapter", imports[3].ID)
		assert.Equal(t, "httpadapter", imports[3].Name)
		assert.Equal(t, `"example.com/test/infrastructure/api/http"`, imports[3].Path)
	}
}

func assertExpectedServices(t *testing.T, container *digen.ContainerDefinition) {
	if assert.Len(t, container.Services, 4) {
		assert.Equal(t, "configuration", container.Services[0].Name)
		assert.Equal(t, "config", container.Services[0].Type.Package)
		assert.Equal(t, "Configuration", container.Services[0].Type.Name)
		assert.False(t, container.Services[0].Type.IsPointer)
		assert.False(t, container.Services[0].HasSetter)
		assert.False(t, container.Services[0].HasCloser)
		assert.False(t, container.Services[0].IsRequired)
		assert.False(t, container.Services[0].IsPublic)
		assert.False(t, container.Services[0].IsExternal)

		assert.Equal(t, "entityRepository", container.Services[1].Name)
		assert.Equal(t, "domain", container.Services[1].Type.Package)
		assert.Equal(t, "EntityRepository", container.Services[1].Type.Name)
		assert.False(t, container.Services[1].Type.IsPointer)
		assert.True(t, container.Services[1].HasSetter)
		assert.True(t, container.Services[1].HasCloser)
		assert.True(t, container.Services[1].IsRequired)
		assert.True(t, container.Services[1].IsPublic)
		assert.True(t, container.Services[1].IsExternal)

		assert.Equal(t, "useCase", container.Services[2].Name)
		assert.Equal(t, "usecase", container.Services[2].Type.Package)
		assert.Equal(t, "FindEntity", container.Services[2].Type.Name)
		assert.True(t, container.Services[2].Type.IsPointer)
		assert.False(t, container.Services[2].HasSetter)
		assert.False(t, container.Services[2].HasCloser)
		assert.False(t, container.Services[2].IsRequired)
		assert.False(t, container.Services[2].IsPublic)
		assert.False(t, container.Services[2].IsExternal)

		assert.Equal(t, "handler", container.Services[3].Name)
		assert.Equal(t, "httpadapter", container.Services[3].Type.Package)
		assert.Equal(t, "GetEntityHandler", container.Services[3].Type.Name)
		assert.True(t, container.Services[3].Type.IsPointer)
		assert.False(t, container.Services[3].HasSetter)
		assert.False(t, container.Services[3].HasCloser)
		assert.False(t, container.Services[3].IsRequired)
		assert.False(t, container.Services[3].IsPublic)
		assert.False(t, container.Services[3].IsExternal)
	}
}
