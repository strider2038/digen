package di_test

import (
	_ "embed"
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
					"config": {Path: "example.com/test/di/config"},
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
			name: "single container with basic types",
			container: &di.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*di.ImportDefinition{
					"time": {Path: "time"},
					"url":  {Path: "net/url"},
				},
				Services: []*di.ServiceDefinition{
					{
						Name:     "StringOption",
						Type:     di.TypeDefinition{Name: "string"},
						IsPublic: true,
					},
					{
						Name:     "StringPointer",
						Type:     di.TypeDefinition{IsPointer: true, Name: "string"},
						IsPublic: true,
					},
					{
						Name:     "IntOption",
						Type:     di.TypeDefinition{Name: "int"},
						IsPublic: true,
					},
					{
						Name:     "TimeOption",
						Type:     di.TypeDefinition{Package: "time", Name: "Time"},
						IsPublic: true,
					},
					{
						Name:     "DurationOption",
						Type:     di.TypeDefinition{Package: "time", Name: "Duration"},
						IsPublic: true,
					},
					{
						Name: "URLOption",
						Type: di.TypeDefinition{Package: "url", Name: "URL"},
					},
					{
						Name: "IntSlice",
						Type: di.TypeDefinition{Name: "int", IsSlice: true},
					},
					{
						Name: "StringMap",
						Type: di.TypeDefinition{Name: "string", Key: &di.TypeDefinition{Name: "string"}},
					},
				},
			},
			assert: func(t *testing.T, files []*di.File) {
				t.Helper()

				require.GreaterOrEqual(t, len(files), 1)
				assert.Equal(t, singleContainerWithBasicTypes, string(files[0].Content))
				assert.Equal(t, singleContainerWithBasicTypesPublicContainer, string(files[2].Content))
			},
		},
		{
			name: "single container with closer",
			container: &di.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*di.ImportDefinition{
					"sql": {Path: "example.com/test/sql"},
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
					"domain": {Path: "example.com/test/domain"},
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
					"domain": {Path: "example.com/test/domain"},
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
								Prefix: "InternalContainerType",
								Name:   "firstService",
								Type: di.TypeDefinition{
									IsPointer: true,
									Package:   "domain",
									Name:      "Service",
								},
								IsPublic: true,
							},
							{
								Prefix: "InternalContainerType",
								Name:   "secondService",
								Type: di.TypeDefinition{
									IsPointer: true,
									Package:   "domain",
									Name:      "Service",
								},
								HasSetter: true,
								HasCloser: true,
							},
							{
								Prefix: "InternalContainerType",
								Name:   "requiredService",
								Type: di.TypeDefinition{
									IsPointer: true,
									Package:   "domain",
									Name:      "Service",
								},
								IsRequired: true,
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
			name: "import alias generation check",
			container: &di.RootContainerDefinition{
				Name: "Container",
				Imports: map[string]*di.ImportDefinition{
					"httpadapter": {Name: "httpadapter", Path: "example.com/test/infrastructure/api/http"},
				},
				Services: []*di.ServiceDefinition{
					{
						Name: "serviceName",
						Type: di.TypeDefinition{
							IsPointer: true,
							Package:   "httpadapter",
							Name:      "ServiceHandler",
						},
						IsPublic: true,
					},
				},
			},
			assert: func(t *testing.T, files []*di.File) {
				t.Helper()

				require.Len(t, files, 3)
				assert.Equal(t, importAliasInternalContainer, string(files[0].Content))
				assert.Equal(t, importAliasLookupContainer, string(files[1].Content))
				assert.Equal(t, importAliasPublicContainer, string(files[2].Content))
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

				require.GreaterOrEqual(t, len(files), 3)
				assert.Contains(
					t, string(files[2].Content),
					`func (c *Container) APIRouter(ctx context.Context) (http.Handler, error) {`,
				)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id := 0
			for _, service := range test.container.Services {
				service.ID = id
				id++
			}
			for _, container := range test.container.Containers {
				for _, service := range container.Services {
					service.ID = id
					id++
				}
			}

			files, err := di.GenerateFiles(test.container, di.GenerationParameters{
				RootPackage: "example.com/test/di",
			})

			require.NoError(t, err)
			test.assert(t, files)
		})
	}
}

var (
	//go:embed testdata/generation/single_container_with_getters_only_internal_container.txt
	singleContainerWithGettersOnlyInternalContainer string
	//go:embed testdata/generation/single_container_with_getters_only_definition_contracts.txt
	singleContainerWithGettersOnlyDefinitionContracts string
	//go:embed testdata/generation/single_container_with_getters_only_public_file.txt
	singleContainerWithGettersOnlyPublicFile string

	//go:embed testdata/generation/single_container_with_service_setter_internal_container.txt
	singleContainerWithServiceSetterInternalContainer string
	//go:embed testdata/generation/single_container_with_service_setter_public_container.txt
	singleContainerWithServiceSetterPublicContainer string

	//go:embed testdata/generation/single_container_with_required_service_internal_container.txt
	singleContainerWithRequiredServiceInternalContainer string
	//go:embed testdata/generation/public_container_with_requirement_file.txt
	publicContainerWithRequirementFile string

	//go:embed testdata/generation/single_container_with_basic_types.txt
	singleContainerWithBasicTypes string
	//go:embed testdata/generation/single_container_with_basic_types_public_container.txt
	singleContainerWithBasicTypesPublicContainer string

	//go:embed testdata/generation/single_container_with_external_service_internal_container.txt
	singleContainerWithExternalServiceInternalContainer string

	//go:embed testdata/generation/single_container_with_static_type_internal_container.txt
	singleContainerWithStaticTypeInternalContainer string

	//go:embed testdata/generation/single_container_with_closer_internal_container.txt
	singleContainerWithCloserInternalContainer string

	//go:embed testdata/generation/two_services_from_one_package_internal_container.txt
	twoServicesFromOnePackageInternalContainer string

	//go:embed testdata/generation/separate_container_internal_container.txt
	separateContainerInternalContainer string
	//go:embed testdata/generation/separate_lookup_container_file.txt
	separateLookupContainerFile string
	//go:embed testdata/generation/separate_container_public_file.txt
	separateContainerPublicFile string

	//go:embed testdata/generation/import_alias_internal_container.txt
	importAliasInternalContainer string
	//go:embed testdata/generation/import_alias_lookup_container.txt
	importAliasLookupContainer string
	//go:embed testdata/generation/import_alias_public_container.txt
	importAliasPublicContainer string
)
