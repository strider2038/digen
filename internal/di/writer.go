package di

import (
	"os"
	"path/filepath"

	"github.com/muonsoft/errors"
	"github.com/spf13/afero"
)

type Writer struct {
	FS        afero.Fs
	WorkDir   string
	Overwrite bool
	Append    bool
}

func NewWriter(fs afero.Fs, workDir string) *Writer {
	return &Writer{FS: fs, WorkDir: workDir}
}

func (w *Writer) WriteFile(file *File) error {
	filename := w.WorkDir + "/"
	if packageDirs[file.Package] != "" {
		filename += packageDirs[file.Package] + "/"
	}
	filename += file.Name

	if isFileExist(w.FS, filename) {
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
	err := w.FS.MkdirAll(dir, 0775)
	if err != nil {
		return errors.Errorf("create dir %s: %w", dir, err)
	}

	err = afero.WriteFile(w.FS, filename, file.Content, 0644)
	if err != nil {
		return errors.Errorf("write file %s: %w", file.Name, err)
	}

	return nil
}

func (w *Writer) append(file *File, filename string) error {
	return appendFile(w.FS, filename, file.Content, 0644)
}

func isFileExist(fs afero.Fs, filename string) bool {
	_, err := fs.Stat(filename)

	return err == nil
}

func appendFile(fs afero.Fs, name string, data []byte, perm os.FileMode) error {
	f, err := fs.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
