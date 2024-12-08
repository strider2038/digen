package di_test

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strider2038/digen/internal/di"
)

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name        string
		testedFiles []string
	}{
		{
			name:        "single container with getters only",
			testedFiles: append(defaultTestedFiles(), "internal/factories/container.go"),
		},
		{name: "single container with service setter"},
		{name: "single container with required service"},
		{name: "single container with static type"},
		{name: "single container with basic types"},
		{name: "single container with closer"},
		{name: "multiple containers"},
		{name: "import alias generation"},
		{name: "override service public name"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			afs := afero.NewMemMapFs()
			setupDefinitionsFile(t, afs, test.name)
			testedFiles := test.testedFiles
			if testedFiles == nil {
				testedFiles = defaultTestedFiles()
			}

			generator := &di.Generator{
				BaseDir:    "di",
				ModulePath: "example.com/test",
				FS:         afs,
				Params: di.GenerationParameters{
					Version: "vX.X.X",
				},
			}
			err := generator.Generate()

			require.NoError(t, err)
			if needToDump() {
				dumpGeneratedFiles(t, afs, test.name, testedFiles)
			}
			assertGeneratedFiles(t, afs, test.name, testedFiles)
		})
	}
}

func defaultTestedFiles() []string {
	return []string{
		"container.go",
		"internal/container.go",
		"internal/lookup/container.go",
	}
}

func setupDefinitionsFile(t *testing.T, fs afero.Fs, testCase string) {
	input := "./testdata/input/" + strcase.ToSnake(testCase) + "_definitions.txt"
	data, err := os.ReadFile(input)
	require.NoError(t, err, "read definitions file")
	err = afero.WriteFile(fs, "./di/internal/definitions/container.go", data, 0644)
	require.NoError(t, err, "write definitions file")
}

func assertGeneratedFiles(t *testing.T, afs afero.Fs, testCase string, testedFiles []string) {
	for _, filename := range testedFiles {
		got, err := afero.ReadFile(afs, "di/"+filename)
		require.NoError(t, err, "read generated file %q", filename)
		want, err := os.ReadFile("./testdata/output/" + formatOutputFilename(testCase, filename) + ".txt")
		require.NoError(t, err, "read expected file %q", filename)
		assert.Equal(t, string(want), string(got), "files not equal for %q", filename)
	}
}

func needToDump() bool {
	v, _ := strconv.ParseBool(os.Getenv("DUMP"))

	return v
}

func dumpGeneratedFiles(t *testing.T, afs afero.Fs, testCase string, testedFiles []string) {
	for _, filename := range testedFiles {
		data, err := afero.ReadFile(afs, "di/"+filename)
		require.NoError(t, err, "read generated file %q", filename)
		err = os.WriteFile("./testdata/output/"+formatOutputFilename(testCase, filename)+".txt", data, 0644)
		require.NoError(t, err, "write generated file %q", filename)
		t.Log("dump file:", filename)
	}
}

func formatOutputFilename(testCase string, filename string) string {
	filename = strings.ReplaceAll(filename, "/", "_")

	return strcase.ToSnake(testCase + "_" + filename)
}
