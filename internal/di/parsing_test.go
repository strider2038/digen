package di_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strider2038/digen/internal/di"
)

var (
	//go:embed testdata/parsing/definitions_container.txt
	definitionsComplexExample string
	//go:embed testdata/parsing/definitions_basic_types.txt
	definitionsBasicTypes string

	//go:embed testdata/parsing/factories.txt
	testFactorySource string
)

func TestParseContainerFromSource(t *testing.T) {
	tests := []struct {
		name   string
		source string
		assert func(t *testing.T, container *di.RootContainerDefinition)
	}{
		{
			name:   "complex example",
			source: definitionsComplexExample,
			assert: func(t *testing.T, container *di.RootContainerDefinition) {
				t.Helper()
				assert.Equal(t, "Container", container.Name)
				assert.Equal(t, "di", container.Package)
				assertExpectedContainerImports(t, container.Imports)
				assertExpectedContainerServices(t, container.Services)
				assertExpectedInternalContainers(t, container.Containers)
			},
		},
		{
			name:   "basic types",
			source: definitionsBasicTypes,
			assert: func(t *testing.T, container *di.RootContainerDefinition) {
				t.Helper()
				assertExpectedBasicTypes(t, container.Services)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			container, err := di.ParseContainerFromSource(test.source)

			require.NoError(t, err)
			require.NotNil(t, container)
			test.assert(t, container)
		})
	}
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
	require.Len(t, services, 4)

	assert.Equal(t, "Configuration", services[0].Name)
	assert.Equal(t, "config", services[0].Type.Package)
	assert.Equal(t, "Configuration", services[0].Type.Name)
	assert.False(t, services[0].Type.IsPointer)
	assert.False(t, services[0].HasSetter)
	assert.False(t, services[0].HasCloser)
	assert.False(t, services[0].IsRequired)
	assert.False(t, services[0].IsPublic)
	assert.False(t, services[0].IsExternal)

	assert.Equal(t, "EntityRepository", services[1].Name)
	assert.Equal(t, "domain", services[1].Type.Package)
	assert.Equal(t, "EntityRepository", services[1].Type.Name)
	assert.False(t, services[1].Type.IsPointer)
	assert.True(t, services[1].HasSetter)
	assert.True(t, services[1].HasCloser)
	assert.True(t, services[1].IsRequired)
	assert.True(t, services[1].IsPublic)
	assert.True(t, services[1].IsExternal)

	assert.Equal(t, "Handler", services[2].Name)
	assert.Equal(t, "httpadapter", services[2].Type.Package)
	assert.Equal(t, "GetEntityHandler", services[2].Type.Name)
	assert.Equal(t, "http_handler.go", services[2].FactoryFileName)
	assert.True(t, services[2].Type.IsPointer)
	assert.False(t, services[2].HasSetter)
	assert.False(t, services[2].HasCloser)
	assert.False(t, services[2].IsRequired)
	assert.False(t, services[2].IsPublic)
	assert.False(t, services[2].IsExternal)

	assert.Equal(t, "Router", services[3].Name)
	assert.Equal(t, "APIRouter", services[3].PublicName)
	assert.Equal(t, "http", services[3].Type.Package)
	assert.Equal(t, "Handler", services[3].Type.Name)
	assert.Equal(t, "", services[3].FactoryFileName)
	assert.False(t, services[3].Type.IsPointer)
	assert.False(t, services[3].HasSetter)
	assert.False(t, services[3].HasCloser)
	assert.False(t, services[3].IsRequired)
	assert.False(t, services[3].IsPublic)
	assert.False(t, services[3].IsExternal)
}

func assertExpectedBasicTypes(t *testing.T, services []*di.ServiceDefinition) {
	require.Len(t, services, 4)

	assert.Equal(t, "StringOption", services[0].Name)
	assert.True(t, services[0].Type.IsBasicType())
	assert.False(t, services[0].Type.IsPointer)
	assert.Equal(t, "", services[0].Type.Package)
	assert.Equal(t, "string", services[0].Type.Name)

	assert.Equal(t, "IntOption", services[1].Name)
	assert.True(t, services[1].Type.IsBasicType())
	assert.False(t, services[1].Type.IsPointer)
	assert.Equal(t, "", services[1].Type.Package)
	assert.Equal(t, "int", services[1].Type.Name)

	assert.Equal(t, "DurationOption", services[2].Name)
	assert.True(t, services[2].Type.IsDuration())
	assert.False(t, services[2].Type.IsPointer)
	assert.Equal(t, "time", services[2].Type.Package)
	assert.Equal(t, "Duration", services[2].Type.Name)

	assert.Equal(t, "StringPointer", services[3].Name)
	assert.True(t, services[3].Type.IsBasicType())
	assert.True(t, services[3].Type.IsPointer)
	assert.Equal(t, "", services[3].Type.Package)
	assert.Equal(t, "string", services[3].Type.Name)
}

func assertExpectedInternalContainers(t *testing.T, containers []*di.ContainerDefinition) {
	if !assert.Len(t, containers, 1) {
		return
	}

	services := containers[0].Services
	if !assert.Len(t, services, 1) {
		return
	}

	assert.Equal(t, "FindEntity", services[0].Name)
	assert.Equal(t, "usecase", services[0].Type.Package)
	assert.Equal(t, "FindEntity", services[0].Type.Name)
	assert.Equal(t, "UseCase", services[0].Prefix)
	assert.True(t, services[0].Type.IsPointer)
	assert.False(t, services[0].HasSetter)
	assert.False(t, services[0].HasCloser)
	assert.False(t, services[0].IsRequired)
	assert.False(t, services[0].IsPublic)
	assert.False(t, services[0].IsExternal)
}
