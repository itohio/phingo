package repository

import (
	"errors"
	"io"
	"io/fs"

	files "github.com/itohio/phingo/pkg/repository/files"
	"github.com/itohio/phingo/pkg/types"
)

type Repository interface {
	io.Closer
	Read() error
	Write() error
	Config() *types.Config

	// Accounts returns a list of accounts that maches the Id.
	// The ID can be either an ID or a Name.
	Accounts(id ...string) []*types.Account
	AccountsCount(id ...string) int

	// Clients returns a list of clients that maches the Id.
	// The ID can be either an ID or a Name.
	Clients(id ...string) []*types.Client
	ClientsCount(id ...string) int

	// Projects returns a list of projects that maches the Id.
	// The ID can be either an ID or a Name.
	Projects(id ...string) []*types.Project
	ProjectsCount(id ...string) int

	Invoices(id ...string) []*types.Invoice
	InvoicesCount(id ...string) int

	Templates(id ...string) []*types.Template
	TemplatesCount(id ...string) int

	SetConfig(*types.Config) error
	SetAccount(*types.Account) error
	SetClient(*types.Client) error
	SetProject(*types.Project) error
	SetInvoice(*types.Invoice) error
	SetTemplate(*types.Template) error

	DelAccount(*types.Account) error
	DelClient(*types.Client) error
	DelProject(*types.Project) error
	DelInvoice(*types.Invoice) error
	DelTemplate(*types.Template) error

	FS() fs.FS
}

// New returns a new repository
func New(url string) (rep Repository, err error) {

	switch {
	case files.CanAccept(url):
		rep, err = files.New(url)
	default:
		return nil, errors.New("repository not found")
	}
	if err == nil {
		err = rep.Read()
	}

	return rep, err
}

// Init initializes a repository from scratch.
func Init(url string) error {
	switch {
	case files.CanAccept(url):
		return files.Init(url)
	default:
		return errors.New("repository not found")
	}
}

func Migrate(url string) error {
	switch {
	case files.CanAccept(url):
		return files.Migrate(url)
	default:
		return errors.New("repository not found")
	}
}
