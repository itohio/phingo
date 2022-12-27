package repository

import (
	"errors"

	"github.com/itohio/phingo/pkg/types"
)

func (r *repository) SetConfig(cfg *types.Config) error {
	r.config = cfg
	r.configModified = true
	return nil
}

func (r *repository) SetAccount(acc *types.Account) error {
	r.accountsModified = true
	if acc.Id == "" {
		acc.Id = acc.MakeId("")
	}
	for i, val := range r.accounts.Accounts {
		if val == acc {
			return nil
		}
		if val.Id == acc.Id {
			r.accounts.Accounts[i] = acc
			return nil
		}
	}
	r.accounts.Accounts = append(r.accounts.Accounts, acc)
	return nil
}

func (r *repository) SetClient(cl *types.Client) error {
	r.clientsModified = true
	if cl.Id == "" {
		cl.Id = cl.MakeId("")
	}
	for i, val := range r.clients.Clients {
		if val == cl {
			return nil
		}
		if val.Id == cl.Id {
			r.clients.Clients[i] = cl
			return nil
		}
	}
	r.clients.Clients = append(r.clients.Clients, cl)
	return nil
}

func (r *repository) SetProject(prj *types.Project) error {
	if prj.Id == "" {
		prj.Id = prj.MakeId()
	}
	r.projectsModified[prj.Id] = modifyStruct{}
	for i, val := range r.projects {
		if val == prj {
			return nil
		}
		if val.Id == prj.Id {
			r.projects[i] = prj
			return nil
		}
	}
	r.projects = append(r.projects, prj)
	return nil
}

func (r *repository) SetInvoice(inv *types.Invoice) error {
	if inv.Id == "" {
		return errors.New("id must be already set")
	}
	r.invoicesModified[inv.Id] = modifyStruct{}
	for i, val := range r.invoices {
		if val == inv {
			return nil
		}
		if val.Id == inv.Id {
			r.invoices[i] = inv
			return nil
		}
	}
	r.invoices = append(r.invoices, inv)
	return nil
}

func (r *repository) SetTemplate(tpl *types.Template) error {
	if tpl.Id == "" {
		tpl.Id = tpl.MakeId()
	}
	r.templatesModified[tpl.Id] = modifyStruct{}
	for i, val := range r.templates {
		if val == tpl {
			return nil
		}
		if val.Id == tpl.Id {
			r.templates[i] = tpl
			return nil
		}
	}
	r.templates = append(r.templates, tpl)
	return nil
}
