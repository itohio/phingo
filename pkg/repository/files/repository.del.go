package repository

import (
	"errors"

	"github.com/itohio/phingo/pkg/types"
)

type IdGetter interface {
	GetId() string
}

func removeId[T IdGetter](arr []T, id string) ([]T, T, error) {
	var tmp T
	N := len(arr)
	for i, val := range arr {
		if val.GetId() == id {
			if i == N-1 {
				return arr[:N-1], arr[N-1], nil
			}
			tmp, arr[i] = arr[i], arr[N-1]
			return arr[:N-1], tmp, nil
		}
	}
	return arr, tmp, errors.New("item not found")
}

func (r *repository) DelAccount(acc *types.Account) error {
	if acc.Id == "" {
		return errors.New("no such account")
	}
	for _, val := range r.projects {
		if val.Account != nil && val.Account.Id == acc.Id {
			return errors.New("account is being used by a project")
		}
	}
	for _, val := range r.invoices {
		if val.Account != nil && val.Account.Id == acc.Id {
			return errors.New("account is being used by an invoice")
		}
	}
	var err error
	r.accounts.Accounts, _, err = removeId(r.accounts.Accounts, acc.Id)
	return err
}

func (r *repository) DelClient(cl *types.Client) error {
	if cl.Id == "" {
		return errors.New("no such client")
	}
	for _, val := range r.projects {
		if val.Client != nil && val.Client.Id == cl.Id {
			return errors.New("client is being used by a project")
		}
	}
	var err error
	r.clients.Clients, _, err = removeId(r.clients.Clients, cl.Id)
	return err
}

func (r *repository) DelProject(prj *types.Project) error {
	if prj.Id == "" {
		return errors.New("no such project")
	}
	for _, val := range r.invoices {
		if val.Project != nil && val.Project.Id == prj.Id {
			return errors.New("project is being used by an invoice")
		}
	}
	var err error
	var prev *types.Project
	r.projects, prev, err = removeId(r.projects, prj.Id)
	if prev != nil {
		r.projectsModified[prev.Id] = modifyStruct{
			delete:   true,
			filename: prev.FileName,
		}
	}
	return err
}

func (r *repository) DelInvoice(inv *types.Invoice) error {
	if inv.Id == "" {
		return errors.New("no such invoice")
	}
	var err error
	var prev *types.Invoice
	r.invoices, prev, err = removeId(r.invoices, inv.Id)
	if prev != nil {
		r.invoicesModified[prev.Id] = modifyStruct{
			delete:   true,
			filename: prev.FileName,
		}
	}
	return err
}

func (r *repository) DelTemplate(tpl *types.Template) error {
	if tpl.Id == "" {
		return errors.New("no such template")
	}
	var err error
	var prev *types.Template
	r.templates, prev, err = removeId(r.templates, tpl.Id)
	if prev != nil {
		r.templatesModified[prev.Id] = modifyStruct{
			delete:   true,
			filename: prev.FileName,
		}
	}
	return err
}
