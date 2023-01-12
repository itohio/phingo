package repository

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"path"
	"strconv"
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
	buf, err := r.readFile(defaultRepo.PathConfigYaml)
	if err != nil {
		return err
	}
	buf, err = yaml2json(buf)
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
	buf, err := r.readFile(defaultRepo.PathAccountsYaml)
	if err != nil {
		return err
	}
	buf, err = yaml2json(buf)
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
	buf, err := r.readFile(defaultRepo.PathClientsYaml)
	if err != nil {
		return err
	}
	buf, err = yaml2json(buf)
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
			if !strings.HasSuffix(pth, ".yaml") {
				return nil
			}
			buf, err := r.readFile(pth)
			if err != nil {
				return err
			}
			buf, err = yaml2json(buf)
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
			if prj.Name == "" {
				return errors.New("badly formatted project: Name")
			}
			if prj.Id == "" {
				return errors.New("badly formatted project: Id")
			}
			prj.FileName = pth
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
			if !strings.HasSuffix(pth, ".yaml") {
				return nil
			}
			b := path.Dir(pth)
			if path.Dir(b) != defaultRepo.PathInvoices {
				return nil
			}
			year, err := strconv.ParseInt(path.Base(b), 10, 32)
			if err != nil {
				return err
			}
			if year < 1971 || year > 3000 {
				return errors.New("invalid year")
			}

			buf, err := r.readFile(pth)
			if err != nil {
				return err
			}
			buf, err = yaml2json(buf)
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
			if _, err := types.ParseTime(inv.IssueDate); err != nil {
				return fmt.Errorf("badly formatted invoice: IssueDate %v", err)
			}
			if inv.Year() != int(year) {
				return errors.New("badly formatted invoice: Year")
			}
			if inv.Id == "" {
				return errors.New("badly formatted invoice: Id")
			}
			if inv.Code == "" {
				return errors.New("badly formatted invoice: Code")
			}
			if inv.Account == nil {
				return errors.New("badly formatted invoice: Account")
			}
			if inv.Client == nil {
				return errors.New("badly formatted invoice: Client")
			}
			inv.Project = r.resolveProject(inv.Project)
			inv.FileName = pth
			r.invoices = append(r.invoices, &inv)
			return nil
		},
	)
}
