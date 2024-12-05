package di_test

import (
	_ "embed"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strider2038/digen/internal/di"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name      string
		container *di.RootContainerDefinition
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
		},
		{
			name: "override service public name",
			container: &di.RootContainerDefinition{
				Name:    "Container",
				Package: "testpkg",
				Imports: map[string]*di.ImportDefinition{
					"http": {Path: "net/http"},
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
				RootPackage:   "example.com/test/di",
				ErrorHandling: di.ErrorHandling{}.Defaults(),
			})

			if needToDump() {
				dumpFiles(t, test.name, files)
			}
			require.NoError(t, err)
			assertFiles(t, test.name, files)
		})
	}
}

func assertFiles(tb testing.TB, testCase string, files []*di.File) {
	for _, file := range files {
		filename := formatFilename(testCase, file)
		content, err := os.ReadFile(filename)
		if err != nil {
			tb.Fatal("read file: ", err)
		}

		assert.Equal(tb, string(content), string(file.Content))
	}
}

func formatFilename(testCase string, file *di.File) string {
	path := strings.ReplaceAll(file.Path(), "/", "_")

	return "./testdata/output/" + strcase.ToSnake(testCase+"_"+path) + ".txt"
}

func needToDump() bool {
	v, _ := strconv.ParseBool(os.Getenv("DUMP"))

	return v
}

func dumpFiles(tb testing.TB, testCase string, files []*di.File) {
	for _, file := range files {
		filename := formatFilename(testCase, file)
		if err := os.WriteFile(filename, file.Content, 0644); err != nil {
			tb.Fatal("write file: ", err)
		}
	}
}
