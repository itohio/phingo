package repository

import (
	"fmt"
	"io"
	"io/fs"
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

func (r *repository) FS() fs.FS {
	return r.fs
}

func (r *repository) Read() error {
	if closer, ok := r.fs.(io.Closer); ok {
		defer closer.Close()
	}
	if err := r.readConfig(); err != nil {
		return fmt.Errorf("readConfig: %v", err)
	}
	if err := r.readAccounts(); err != nil {
		return fmt.Errorf("readAccounts: %v", err)
	}
	if err := r.readClients(); err != nil {
		return fmt.Errorf("readClients: %v", err)
	}
	if err := r.readTemplates(); err != nil {
		return fmt.Errorf("readTemplates: %v", err)
	}
	if err := r.readProjects(); err != nil {
		return fmt.Errorf("readProjects: %v", err)
	}
	if err := r.readInvoices(); err != nil {
		return fmt.Errorf("readInvoices: %v", err)
	}
	return nil
}

func (r *repository) Write() error {
	if closer, ok := r.fs.(io.Closer); ok {
		defer closer.Close()
	}
	if err := r.writeConfig(); err != nil {
		return fmt.Errorf("writeConfig: %v", err)
	}
	if err := r.writeAccounts(); err != nil {
		return fmt.Errorf("writeAccounts: %v", err)
	}
	if err := r.writeClients(); err != nil {
		return fmt.Errorf("writeClients: %v", err)
	}
	if err := r.writeTemplates(); err != nil {
		return fmt.Errorf("writeTemplates: %v", err)
	}
	if err := r.writeProjects(); err != nil {
		return fmt.Errorf("writeProjects: %v", err)
	}
	if err := r.writeInvoices(); err != nil {
		return fmt.Errorf("writeInvoices: %v", err)
	}
	return nil
}

func (r *repository) Config() *types.Config {
	return r.config
}

func accountsPredicate(mid map[string]struct{}) func(a *types.Account) bool {
	return func(a *types.Account) bool {
		if len(mid) == 0 {
			return true
		}
		if _, ok := mid[a.Id]; ok {
			return true
		}
		if _, ok := mid["name:"+a.Name]; ok {
			return true
		}
		return false
	}
}

func makeIdMap(id []string) map[string]struct{} {
	mid := make(map[string]struct{}, len(id))
	for _, id := range id {
		if id == "" {
			continue
		}
		mid[id] = struct{}{}
	}
	return mid
}

func (r *repository) Accounts(ids ...string) []*types.Account {
	mid := makeIdMap(ids)
	return types.Filter(r.accounts.Accounts, accountsPredicate(mid))
}

func (r *repository) AccountsCount(ids ...string) int {
	mid := makeIdMap(ids)

	return types.Count(r.accounts.Accounts, accountsPredicate(mid))
}

func clientsPredicate(mid map[string]struct{}) func(a *types.Client) bool {
	return func(a *types.Client) bool {
		if len(mid) == 0 {
			return true
		}
		if _, ok := mid[a.Id]; ok {
			return true
		}
		if _, ok := mid["name:"+a.Name]; ok {
			return true
		}
		return false
	}
}

func (r *repository) Clients(ids ...string) []*types.Client {
	mid := makeIdMap(ids)

	return types.Filter(r.clients.Clients, clientsPredicate(mid))
}

func (r *repository) ClientsCount(ids ...string) int {
	mid := makeIdMap(ids)

	return types.Count(r.clients.Clients, clientsPredicate(mid))
}

func projectsPredicate(mid map[string]struct{}) func(a *types.Project) bool {
	return func(a *types.Project) bool {
		if len(mid) == 0 {
			return true
		}
		if _, ok := mid[a.Id]; ok {
			return true
		}
		if _, ok := mid["name:"+a.Name]; ok {
			return true
		}
		return false
	}
}

func (r *repository) Projects(ids ...string) []*types.Project {
	mid := makeIdMap(ids)

	return types.Filter(r.projects, projectsPredicate(mid))
}

func (r *repository) ProjectsCount(ids ...string) int {
	mid := makeIdMap(ids)

	return types.Count(r.projects, projectsPredicate(mid))
}

func invoicesPredicate(mid map[string]struct{}) func(*types.Invoice) bool {
	return func(a *types.Invoice) bool {
		if len(mid) == 0 {
			return true
		}
		if _, ok := mid[a.Id]; ok {
			return true
		}
		year := fmt.Sprintf("year:%d", a.Year())
		if _, ok := mid[year]; ok {
			return true
		}
		if a.Client == nil {
			return false
		}
		client := fmt.Sprintf("client:%s", a.Client.Name)
		if _, ok := mid[client]; ok {
			return true
		}
		if _, ok := mid[fmt.Sprintf("%s;%s", year, client)]; ok {
			return true
		}
		if a.Project == nil {
			return false
		}
		project := fmt.Sprintf("project:%s", a.Project.Name)
		if _, ok := mid[project]; ok {
			return true
		}
		if _, ok := mid[fmt.Sprintf("%s;%s", year, project)]; ok {
			return true
		}
		if _, ok := mid[fmt.Sprintf("%s;%s", client, project)]; ok {
			return true
		}
		if _, ok := mid[fmt.Sprintf("%s;%s;%s", year, client, project)]; ok {
			return true
		}
		return false
	}
}

func (r *repository) Invoices(ids ...string) []*types.Invoice {
	mid := makeIdMap(ids)

	return types.Filter(r.invoices, invoicesPredicate(mid))
}

func (r *repository) InvoicesCount(ids ...string) int {
	mid := makeIdMap(ids)

	return types.Count(r.invoices, invoicesPredicate(mid))
}

func templatesPredicate(mid map[string]struct{}) func(a *types.Template) bool {
	return func(a *types.Template) bool {
		if len(mid) == 0 {
			return true
		}
		if _, ok := mid[a.Id]; ok {
			return true
		}
		return false
	}
}

func (r *repository) Templates(ids ...string) []*types.Template {
	mid := makeIdMap(ids)

	return types.Filter(r.templates, templatesPredicate(mid))
}

func (r *repository) TemplatesCount(ids ...string) int {
	mid := makeIdMap(ids)

	return types.Count(r.templates, templatesPredicate(mid))

}

func (r *repository) Close() error {
	return nil
}
