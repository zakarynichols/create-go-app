package dir

// While it has logic to make files its really only used for directories.

import (
	"io"
	"os"
	"path/filepath"
)

func CpdirAll(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		file, err := os.Stat(srcPath)
		if err != nil {
			return err
		}

		maskedMode := file.Mode() & os.ModeType

		if maskedMode == os.ModeDir {
			err := mkdir(dstPath, 0755)
			if err != nil {
				return err
			}

			err = CpdirAll(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			_, err := cpFile(srcPath, dstPath) // return the written bytes?
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func cpFile(src, dst string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()
	written, err := io.Copy(dstFile, srcFile)
	return written, err
}

func mkdir(dir string, perm os.FileMode) error {
	err := os.MkdirAll(dir, perm)
	if err != nil {
		return err
	}

	return nil
}
