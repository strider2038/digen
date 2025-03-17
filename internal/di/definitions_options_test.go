package di_test

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/strider2038/digen/internal/di"
)

func TestOptionsParser_ParseServiceDefinitionOptions(t *testing.T) {
	tests := []struct {
		name  string
		field *ast.Field
		want  di.ServiceDefinitionsOptions
	}{
		{
			name:  "tags: public flag",
			field: &ast.Field{Tag: &ast.BasicLit{Value: "`di:\"public\"`"}},
			want: di.ServiceDefinitionsOptions{
				Flags: []string{"public"},
			},
		},
		{
			name:  "tags: set flag",
			field: &ast.Field{Tag: &ast.BasicLit{Value: "`di:\"set\"`"}},
			want: di.ServiceDefinitionsOptions{
				Flags: []string{"set"},
			},
		},
		{
			name:  "tags: close flag",
			field: &ast.Field{Tag: &ast.BasicLit{Value: "`di:\"close\"`"}},
			want: di.ServiceDefinitionsOptions{
				Flags: []string{"close"},
			},
		},
		{
			name:  "tags: required flag",
			field: &ast.Field{Tag: &ast.BasicLit{Value: "`di:\"required\"`"}},
			want: di.ServiceDefinitionsOptions{
				Flags: []string{"required"},
			},
		},
		{
			name:  "tags: multiple flags",
			field: &ast.Field{Tag: &ast.BasicLit{Value: "`di:\"set,close\"`"}},
			want: di.ServiceDefinitionsOptions{
				Flags: []string{"set", "close"},
			},
		},
		{
			name:  "tags: public name",
			field: &ast.Field{Tag: &ast.BasicLit{Value: "`public_name:\"PublicName\"`"}},
			want: di.ServiceDefinitionsOptions{
				PublicName: "PublicName",
			},
		},
		{
			name:  "tags: factory file",
			field: &ast.Field{Tag: &ast.BasicLit{Value: "`factory_file:\"FactoryFile\"`"}},
			want: di.ServiceDefinitionsOptions{
				FactoryFilename: "FactoryFile",
			},
		},
		{
			name:  "tags: factory pkg",
			field: &ast.Field{Tag: &ast.BasicLit{Value: "`factory_pkg:\"FactoryPkg\"`"}},
			want: di.ServiceDefinitionsOptions{
				FactoryPackage: "FactoryPkg",
			},
		},
		{
			name: "comments: public flag",
			field: &ast.Field{
				Doc: &ast.CommentGroup{List: []*ast.Comment{{Text: "// di:public"}}},
			},
			want: di.ServiceDefinitionsOptions{
				Flags: []string{"public"},
			},
		},
		{
			name: "comments: set flag",
			field: &ast.Field{
				Doc: &ast.CommentGroup{List: []*ast.Comment{{Text: "// di:set"}}},
			},
			want: di.ServiceDefinitionsOptions{
				Flags: []string{"set"},
			},
		},
		{
			name: "comments: close flag",
			field: &ast.Field{
				Doc: &ast.CommentGroup{List: []*ast.Comment{{Text: "// di:close"}}},
			},
			want: di.ServiceDefinitionsOptions{
				Flags: []string{"close"},
			},
		},
		{
			name: "comments: required flag",
			field: &ast.Field{
				Doc: &ast.CommentGroup{List: []*ast.Comment{{Text: "// di:required"}}},
			},
			want: di.ServiceDefinitionsOptions{
				Flags: []string{"required"},
			},
		},
		{
			name: "comments: multiple flags",
			field: &ast.Field{
				Doc: &ast.CommentGroup{List: []*ast.Comment{{Text: "// di:set,close"}}},
			},
			want: di.ServiceDefinitionsOptions{
				Flags: []string{"set", "close"},
			},
		},
		{
			name: "comments: public name",
			field: &ast.Field{
				Doc: &ast.CommentGroup{List: []*ast.Comment{{Text: "// di:public_name:PublicName"}}},
			},
			want: di.ServiceDefinitionsOptions{
				PublicName: "PublicName",
			},
		},
		{
			name: "comments: factory file",
			field: &ast.Field{
				Doc: &ast.CommentGroup{List: []*ast.Comment{{Text: "// di:factory_file:FactoryFile"}}},
			},
			want: di.ServiceDefinitionsOptions{
				FactoryFilename: "FactoryFile",
			},
		},
		{
			name: "comments: factory pkg",
			field: &ast.Field{
				Doc: &ast.CommentGroup{List: []*ast.Comment{{Text: "// di:factory_pkg:FactoryPkg"}}},
			},
			want: di.ServiceDefinitionsOptions{
				FactoryPackage: "FactoryPkg",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := di.OptionsParser{Logger: &testingLogger{tb: t}}

			got := parser.ParseServiceDefinitionOptions(test.field)

			assert.Equal(t, test.want, got)
		})
	}
}

type testingLogger struct {
	tb testing.TB
}

func (t *testingLogger) Debug(a ...any)           { t.tb.Log(a...) }
func (t *testingLogger) Info(a ...interface{})    { t.tb.Log(a...) }
func (t *testingLogger) Success(a ...interface{}) { t.tb.Log(a...) }
func (t *testingLogger) Warning(a ...interface{}) { t.tb.Log(a...) }
