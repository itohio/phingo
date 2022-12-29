package repository

import (
	"errors"
	"io/fs"
	"os"
)

const (
	PathConfigYaml   = "config.yaml"
	PathAccountsYaml = "accounts.yaml"
	PathClientsYaml  = "clients.yaml"
	PathProjects     = "projects"
	PathTemplates    = "templates"
	PathInvoices     = "invoices"
)

type RWFS interface {
	fs.FS

	MkDirAll(path string, perm fs.FileMode) error
	Create(path string) (*os.File, error)
	Rename(src, dst string) error
	Remove(path string) error
}

type OSWrapper struct {
	rfs       fs.FS
	_MkDirAll func(path string, perm fs.FileMode) error
	_Create   func(path string) (*os.File, error)
	_Rename   func(src, dst string) error
	_Remove   func(path string) error
}

var _ RWFS = &OSWrapper{}

func NewOSWrapper(fsys fs.FS) (*OSWrapper, error) {
	if _, ok := fsys.(fs.StatFS); !ok {
		return nil, errors.New("must support StatFS")
	}
	if _, ok := fsys.(fs.ReadDirFS); !ok {
		return nil, errors.New("must support ReadDirFS")
	}

	ret := &OSWrapper{
		rfs:       fsys,
		_MkDirAll: os.MkdirAll,
		_Create:   os.Create,
		_Rename:   os.Rename,
		_Remove:   os.Remove,
	}

	return ret, nil
}

func (w *OSWrapper) Open(name string) (fs.File, error)     { return w.rfs.Open(name) }
func (w *OSWrapper) Stat(name string) (fs.FileInfo, error) { return w.rfs.(fs.StatFS).Stat(name) }
func (w *OSWrapper) ReadDir(name string) ([]fs.DirEntry, error) {
	return w.rfs.(fs.ReadDirFS).ReadDir(name)
}

func (w *OSWrapper) MkDirAll(path string, perm fs.FileMode) error { return w._MkDirAll(path, perm) }
func (w *OSWrapper) Create(path string) (*os.File, error)         { return w._Create(path) }
func (w *OSWrapper) Rename(src, dst string) error                 { return w._Rename(src, dst) }
func (w *OSWrapper) Remove(path string) error                     { return w._Remove(path) }
