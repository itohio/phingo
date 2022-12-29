package repository

import (
	"errors"
	"io/fs"
	"os"
)

type tarFS struct {
}

func newTarFS(url string) (*tarFS, error) {
	return nil, errors.New("not implemented")
}

func (w *tarFS) Open(name string) (fs.File, error) {
	return nil, errors.New("not implemented")
}

func (w *tarFS) Stat(name string) (fs.FileInfo, error) {
	return nil, errors.New("not implemented")
}

func (w *tarFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return nil, errors.New("not implemented")
}

func (w *tarFS) MkDirAll(path string, perm fs.FileMode) error {
	return errors.New("not implemented")
}

func (w *tarFS) Create(path string) (*os.File, error) {
	return nil, errors.New("not implemented")
}

func (w *tarFS) Rename(src, dst string) error {
	return errors.New("not implemented")
}

func (w *tarFS) Remove(path string) error {
	return errors.New("not implemented")
}
