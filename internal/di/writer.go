package di

import (
	"os"
	"path/filepath"

	"github.com/muonsoft/errors"
)

type Writer struct {
	WorkDir   string
	Overwrite bool
	Append    bool
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

	if isFileExist(filename) {
		if w.Append {
			return w.append(file, filename)
		}
		if !w.Overwrite {
			return errors.Errorf("cannot write to file %s: %w", filename, ErrFileAlreadyExists)
		}
	}

	return w.write(file, filename)
}

func (w *Writer) write(file *File, filename string) error {
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

func (w *Writer) append(file *File, filename string) error {
	return appendFile(filename, file.Content, 0644)
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}

func appendFile(name string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
