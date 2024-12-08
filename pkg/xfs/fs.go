package xfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

func EnsureDir(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	return os.MkdirAll(dir, perm)
}

func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func Cwd() (string, error) {
	return os.Getwd()
}

func Absolte(path string) (string, error) {
	if !filepath.IsAbs(path) {
		return filepath.Abs(path)
	}

	return path, nil
}

func Resolve(relative string, base string) (string, error) {
	if filepath.IsAbs(relative) {
		return relative, nil
	}

	if relative[0] == '.' && (relative[1] == '/' || relative[1] == '\\') {
		return filepath.Abs(filepath.Join(base, relative[2:]))
	}

	return filepath.Abs(filepath.Join(base, relative))
}

type WalkDirFunc func(path string, d fs.DirEntry, err error) error
