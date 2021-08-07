package digen

import "bytes"

type PackageType int

const (
	unknownPackage PackageType = iota
	PublicPackage
	InternalPackage
	DefinitionsPackage
	lastPackage
)

var packageDirs = [lastPackage]string{
	InternalPackage:    "internal",
	DefinitionsPackage: "internal/definitions",
}

type File struct {
	Package PackageType
	Name    string
	Content []byte
}

type FileBuilder struct {
	fileName    string
	packageName string
	packageType PackageType
	imports     []string
	body        bytes.Buffer
}

func NewFileBuilder(filename, packageName string, packageType PackageType) *FileBuilder {
	return &FileBuilder{
		fileName:    filename,
		packageName: packageName,
		packageType: packageType,
	}
}

func (b *FileBuilder) AddImport(imp string) {
	for _, existingImport := range b.imports {
		if existingImport == imp {
			return
		}
	}

	b.imports = append(b.imports, imp)
}

func (b *FileBuilder) Write(p []byte) (n int, err error) {
	return b.body.Write(p)
}

func (b *FileBuilder) WriteString(s string) (n int, err error) {
	return b.body.WriteString(s)
}

func (b *FileBuilder) GetFile() *File {
	var buffer bytes.Buffer

	buffer.WriteString("package " + b.packageName + "\n\n")

	if len(b.imports) > 0 {
		buffer.WriteString("import (\n")
		for _, imp := range b.imports {
			buffer.WriteString("\t" + imp + "\n")
		}
		buffer.WriteString(")\n")
	}

	b.body.WriteTo(&buffer)

	return &File{
		Package: b.packageType,
		Name:    b.fileName,
		Content: buffer.Bytes(),
	}
}
