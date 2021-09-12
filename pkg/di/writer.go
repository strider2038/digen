package di

import (
	"os"

	"github.com/pkg/errors"
)

type Writer struct {
	WorkDir   string
	Overwrite bool
	Heading   []byte
}

func NewWriter(workDir string) *Writer {
	return &Writer{WorkDir: workDir}
}

func (w *Writer) WriteFile(file *File) error {
	dir := w.WorkDir + "/" + packageDirs[file.Package]
	filename := dir + "/" + file.Name

	if !w.Overwrite && isFileExist(filename) {
		return errors.Wrapf(ErrFileAlreadyExists, "cannot write to file %s", filename)
	}

	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return errors.Wrapf(err, "failed to create dir %s", dir)
	}

	err = os.WriteFile(filename, append(w.Heading, file.Content...), 0644)
	if err != nil {
		return errors.WithMessagef(err, "failed to write %s", file.Name)
	}

	return nil
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}
