package repository

import (
	"io"
	"io/fs"
	"log"
	"path"
	"strings"

	defaultRepo "github.com/itohio/phingo/pkg/repository/default"
	"github.com/itohio/phingo/pkg/types"

	"google.golang.org/protobuf/encoding/protojson"
)

func (r *repository) readFile(pth string) ([]byte, error) {
	f, err := r.fs.Open(pth)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (r *repository) readConfig() error {
	buf, err := r.readFile(defaultRepo.PathConfig)
	if err != nil {
		return err
	}
	r.configModified = false
	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}.Unmarshal(buf, r.config)
}

func (r *repository) readAccounts() error {
	buf, err := r.readFile(defaultRepo.PathAccounts)
	if err != nil {
		return err
	}
	r.accountsModified = false
	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}.Unmarshal(buf, r.accounts)
}

func (r *repository) readClients() error {
	buf, err := r.readFile(defaultRepo.PathClients)
	if err != nil {
		return err
	}
	r.clientsModified = false
	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}.Unmarshal(buf, r.clients)
}

func (r *repository) readTemplates() error {
	r.templates = nil
	r.templatesModified = make(map[string]modifyStruct)

	return fs.WalkDir(
		r.fs,
		defaultRepo.PathTemplates,
		func(pth string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(pth, ".md") {
				return nil
			}
			buf, err := r.readFile(pth)
			if err != nil {
				return err
			}

			id := path.Base(pth)
			id = id[:len(id)-3]
			tpl := &types.Template{
				Id:       id,
				What:     id,
				FileName: pth,
				Text:     buf,
			}
			r.templates = append(r.templates, tpl)
			return nil
		},
	)
}

func (r *repository) readProjects() error {
	r.projects = nil
	r.projectsModified = make(map[string]modifyStruct)

	return fs.WalkDir(
		r.fs,
		defaultRepo.PathProjects,
		func(pth string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(pth, ".json") {
				return nil
			}
			buf, err := r.readFile(pth)
			if err != nil {
				return err
			}
			var prj types.Project
			err = protojson.UnmarshalOptions{
				DiscardUnknown: true,
				AllowPartial:   true,
			}.Unmarshal(buf, &prj)
			if err != nil {
				return err
			}
			prj.FileName = pth
			id := path.Base(pth)
			id = id[:len(id)-5]
			if prj.Id == "" {
				prj.Id = id
			}
			if prj.Name == "" {
				prj.Name = strings.ToTitle(id)
			}
			prj.Account = r.resolveAccount(prj.Account)
			prj.Client = r.resolveClient(prj.Client)
			r.projects = append(r.projects, &prj)
			return nil
		},
	)
}

func (r *repository) resolveAccount(acc *types.Account) *types.Account {
	if acc == nil {
		return nil
	}
	if acc.Id == "" {
		return nil
	}
	for _, a := range r.accounts.Accounts {
		if a.Id == acc.Id {
			return a
		}
	}
	log.Println("resolveAccount: Could not find an account with id=", acc.Id)
	return nil
}

func (r *repository) resolveClient(c *types.Client) *types.Client {
	if c == nil {
		return nil
	}
	if c.Id == "" {
		return nil
	}
	for _, a := range r.clients.Clients {
		if a.Id == c.Id {
			return a
		}
	}
	log.Println("resolveClient: Could not find a client with id=", c.Id)
	return nil
}

func (r *repository) resolveProject(p *types.Project) *types.Project {
	if p == nil {
		return nil
	}
	if p.Id == "" {
		return nil
	}
	for _, a := range r.projects {
		if a.Id == p.Id {
			return a
		}
	}
	log.Println("resolveProject: Could not find a project with id=", p.Id)
	return nil
}

func (r *repository) readInvoices() error {
	r.invoices = nil
	r.invoicesModified = make(map[string]modifyStruct)

	return fs.WalkDir(
		r.fs,
		defaultRepo.PathInvoices,
		func(pth string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(pth, ".json") {
				return nil
			}
			b, id := path.Dir(pth), path.Base(pth)
			if path.Dir(b) != defaultRepo.PathInvoices {
				return nil
			}

			buf, err := r.readFile(pth)
			if err != nil {
				return err
			}
			var inv types.Invoice
			err = protojson.UnmarshalOptions{
				DiscardUnknown: true,
				AllowPartial:   true,
			}.Unmarshal(buf, &inv)
			if err != nil {
				return err
			}
			inv.Year = path.Dir(b)
			inv.FileName = pth
			if inv.Id == "" {
				id = id[:len(id)-5]
				inv.Id = id
			}
			inv.Account = r.resolveAccount(inv.Account)
			inv.Project = r.resolveProject(inv.Project)
			r.invoices = append(r.invoices, &inv)
			return nil
		},
	)
}
