package fsys

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

var EmbedPath string

type fileService interface {
	Create(name string) (*os.File, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	ReadAll(r io.Reader) ([]byte, error)
	Mkdir(name string, perm os.FileMode) error
}

type FileReaderCloser interface {
	Read(p []byte) (n int, err error)
	Close() error
}

type opener interface {
	Open(name string) (FileReaderCloser, error)
}

func Output(name string, path string, isDir bool, o opener, fs fileService) error {
	// Remove the 'embed' string from the path.
	r := strings.Replace(path, EmbedPath, "", -1)

	// Join the new app's directory name to the new string.
	dst := filepath.Join(name, strings.TrimPrefix(r, name))

	// Create directories if they don't exist.
	if isDir {
		err := fs.Mkdir(dst, os.FileMode(0777))
		if os.IsExist(err) {
			return nil
		} else {
			return err
		}
	}

	f, err := o.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := fs.ReadAll(f)
	if err != nil {
		return err
	}

	dstFile, err := fs.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	err = fs.WriteFile(dst, b, os.FileMode(0777))
	if err != nil {
		return err
	}

	return nil
}

type FileDescriptor interface {
	Name() string
	IsDir() bool
}

type FileReaderWriter interface {
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
}

// Change go import paths in all files.
func ReplaceImports(pattern string, path string, old string, new string, fd FileDescriptor, frw FileReaderWriter) error {
	if fd.IsDir() {
		return nil
	}

	matched, err := filepath.Match(pattern, fd.Name())

	if err != nil {
		return err
	}

	if matched {
		read, err := frw.ReadFile(path)
		if err != nil {
			return err
		}

		newContents := strings.Replace(string(read), old, new, -1)

		err = frw.WriteFile(path, []byte(newContents), os.FileMode(0777))

		if err != nil {
			return err
		}
	}

	return nil
}
