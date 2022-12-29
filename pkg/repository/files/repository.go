package repository

import (
	"io"
	"os"
	"strings"

	defaultRepo "github.com/itohio/phingo/pkg/repository/default"
	"github.com/itohio/phingo/pkg/types"
)

type modifyStruct struct {
	delete   bool
	filename string
}

type repository struct {
	url       string
	fs        defaultRepo.RWFS
	config    *types.Config
	accounts  *types.Accounts
	clients   *types.Clients
	templates []*types.Template
	projects  []*types.Project
	invoices  []*types.Invoice

	configModified    bool
	accountsModified  bool
	clientsModified   bool
	templatesModified map[string]modifyStruct
	projectsModified  map[string]modifyStruct
	invoicesModified  map[string]modifyStruct
}

func New(url string) (*repository, error) {
	ret := &repository{
		url:               url,
		config:            defaultRepo.DefaultConfig(),
		accounts:          defaultRepo.DefaultAccounts(),
		clients:           defaultRepo.DefaultClients(),
		templates:         defaultRepo.DefaultTemplates(),
		projects:          defaultRepo.DefaultProjects(),
		invoices:          defaultRepo.DefaultInvoices(),
		templatesModified: make(map[string]modifyStruct),
		projectsModified:  make(map[string]modifyStruct),
		invoicesModified:  make(map[string]modifyStruct),
		configModified:    true,
		accountsModified:  true,
		clientsModified:   true,
	}

	for _, val := range ret.templates {
		ret.templatesModified[val.Id] = modifyStruct{}
	}
	for _, val := range ret.projects {
		ret.projectsModified[val.Id] = modifyStruct{}
	}
	for _, val := range ret.invoices {
		ret.invoicesModified[val.Id] = modifyStruct{}
	}

	switch {
	case strings.HasSuffix(url, "tar.gz"):
		fallthrough
	case strings.HasSuffix(url, "tar"):
		fsys, err := newTarFS(url)
		if err != nil {
			return nil, err
		}
		ret.fs = fsys
	default:
		fsys, err := defaultRepo.NewOSWrapper(os.DirFS(url))
		if err != nil {
			return nil, err
		}
		ret.fs = fsys
	}

	return ret, nil
}

func (r *repository) Read() error {
	if closer, ok := r.fs.(io.Closer); ok {
		defer closer.Close()
	}
	if err := r.readConfig(); err != nil {
		return err
	}
	if err := r.readAccounts(); err != nil {
		return err
	}
	if err := r.readClients(); err != nil {
		return err
	}
	if err := r.readTemplates(); err != nil {
		return err
	}
	if err := r.readProjects(); err != nil {
		return err
	}
	if err := r.readInvoices(); err != nil {
		return err
	}
	return nil
}

func (r *repository) Write() error {
	if closer, ok := r.fs.(io.Closer); ok {
		defer closer.Close()
	}
	if err := r.writeConfig(); err != nil {
		return err
	}
	if err := r.writeAccounts(); err != nil {
		return err
	}
	if err := r.writeClients(); err != nil {
		return err
	}
	if err := r.writeTemplates(); err != nil {
		return err
	}
	if err := r.writeProjects(); err != nil {
		return err
	}
	if err := r.writeInvoices(); err != nil {
		return err
	}
	return nil
}

func (r *repository) Config() *types.Config {
	return r.config
}

func (r *repository) Accounts(id ...string) []*types.Account {
	mid := make(map[string]struct{}, len(id))
	for _, id := range id {
		mid[id] = struct{}{}
	}

	acc := make([]*types.Account, 0, len(mid))
	for _, a := range r.accounts.Accounts {
		if _, ok := mid[a.Id]; !ok && len(id) != 0 {
			if _, ok := mid[a.Name]; !ok && len(id) != 0 {
				continue
			}
		}
		acc = append(acc, a)
	}

	return acc
}

func (r *repository) Clients(id ...string) []*types.Client {
	mid := make(map[string]struct{}, len(id))
	for _, id := range id {
		mid[id] = struct{}{}
	}

	acc := make([]*types.Client, 0, len(mid))
	for _, a := range r.clients.Clients {
		if _, ok := mid[a.Id]; !ok && len(id) != 0 {
			if _, ok := mid[a.Name]; !ok && len(id) != 0 {
				continue
			}
		}
		acc = append(acc, a)
	}

	return acc
}

func (r *repository) Projects(id ...string) []*types.Project {
	mid := make(map[string]struct{}, len(id))
	for _, id := range id {
		mid[id] = struct{}{}
	}

	acc := make([]*types.Project, 0, len(mid))
	for _, a := range r.projects {
		if _, ok := mid[a.Id]; !ok && len(id) != 0 {
			if _, ok := mid[a.Name]; !ok && len(id) != 0 {
				continue
			}
		}
		acc = append(acc, a)
	}

	return acc
}

func (r *repository) Invoices(id ...string) []*types.Invoice {
	mid := make(map[string]struct{}, len(id))
	for _, id := range id {
		mid[id] = struct{}{}
	}

	acc := make([]*types.Invoice, 0, len(mid))
	for _, a := range r.invoices {
		if _, ok := mid[a.Id]; !ok && len(id) != 0 {
			continue
		}
		acc = append(acc, a)
	}

	return acc
}

func (r *repository) Templates(id ...string) []*types.Template {
	mid := make(map[string]struct{}, len(id))
	for _, id := range id {
		mid[id] = struct{}{}
	}

	acc := make([]*types.Template, 0, len(mid))
	for _, a := range r.templates {
		if _, ok := mid[a.Id]; !ok && len(id) != 0 {
			continue
		}
		acc = append(acc, a)
	}

	return acc
}

func (r *repository) Close() error {
	return nil
}
