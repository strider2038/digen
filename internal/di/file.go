package di

import (
	"bytes"
	"go/format"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/muonsoft/errors"
)

type PackageType int

const (
	unknownPackage PackageType = iota
	PublicPackage
	InternalPackage
	DefinitionsPackage
	FactoriesPackage
	LookupPackage
	lastPackage
)

var packageDirs = [lastPackage]string{
	InternalPackage:    "internal",
	DefinitionsPackage: "internal/definitions",
	FactoriesPackage:   "internal/factories",
	LookupPackage:      "internal/lookup",
}

type File struct {
	Package PackageType
	Name    string
	Content []byte
	Append  bool
}

func (f *File) Path() string {
	path := packageDirs[f.Package]
	if path != "" {
		path += "/"
	}

	return path + f.Name
}

func (f *File) IsEmpty() bool {
	return len(f.Content) == 0
}

type FileBuilder struct {
	file        *jen.File
	fileName    string
	packageName string
	packageType PackageType
	body        bytes.Buffer
}

func NewFileBuilder(filename, packageName string, packageType PackageType) *FileBuilder {
	return &FileBuilder{
		file:        jen.NewFile(packageName),
		fileName:    filename,
		packageName: packageName,
		packageType: packageType,
	}
}

func (b *FileBuilder) AddImportAliases(imports map[string]*ImportDefinition) {
	for _, definition := range imports {
		if definition.Name != "" {
			b.file.ImportAlias(strings.Trim(definition.Path, `"`), definition.Name)
		}
	}
}

func (b *FileBuilder) Add(code ...jen.Code) *jen.Statement {
	return b.file.Add(code...)
}

func (b *FileBuilder) GetFile() (*File, error) {
	var buffer bytes.Buffer

	if err := b.file.Render(&buffer); err != nil {
		return nil, errors.Errorf("render %s: %w", b.fileName, err)
	}

	content, err := format.Source(buffer.Bytes())
	if err != nil {
		content = buffer.Bytes()
	}

	return &File{
		Package: b.packageType,
		Name:    b.fileName,
		Content: content,
	}, nil
}
