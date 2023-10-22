package di

import (
	"os"
	"path/filepath"

	"github.com/muonsoft/errors"
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
	filename := w.WorkDir + "/"
	if packageDirs[file.Package] != "" {
		filename += packageDirs[file.Package] + "/"
	}
	filename += file.Name

	if !w.Overwrite && isFileExist(filename) {
		return errors.Errorf("cannot write to file %s: %w", filename, ErrFileAlreadyExists)
	}

	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return errors.Errorf("create dir %s: %w", dir, err)
	}

	err = os.WriteFile(filename, append(w.Heading, file.Content...), 0644)
	if err != nil {
		return errors.Errorf("write file %s: %w", file.Name, err)
	}

	return nil
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}
