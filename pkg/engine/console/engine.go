package engine

import (
	"io"
	"text/template"

	"github.com/itohio/phingo/pkg/types"
)

type engine struct {
	cfg *types.Config
}

func New(config *types.Config) (*engine, error) {
	return &engine{
		cfg: config,
	}, nil
}

func (e *engine) Meta(tpl *types.Template) (map[string]interface{}, error) {
	return nil, nil
}

func (e *engine) ExportAccounts(writer io.Writer, tpl *types.Template, accounts []*types.Account) error {
	context := &types.AccountTemplateContext{
		Config:   e.cfg,
		Template: tpl,
		Accounts: accounts,
	}
	t, err := template.New(tpl.What).Parse(string(tpl.Text))
	if err != nil {
		return nil
	}

	return t.Execute(writer, context)
}

func (e *engine) ExportClients(writer io.Writer, tpl *types.Template, clients []*types.Client) error {
	context := &types.ClientTemplateContext{
		Config:   e.cfg,
		Template: tpl,
		Clients:  clients,
	}
	t, err := template.New(tpl.What).Parse(string(tpl.Text))
	if err != nil {
		return nil
	}

	return t.Execute(writer, context)
}

func (e *engine) ExportProjects(writer io.Writer, tpl *types.Template, projects []*types.Project) error {
	context := &types.ProjectTemplateContext{
		Config:   e.cfg,
		Template: tpl,
		Projects: projects,
	}
	t, err := template.New(tpl.What).Parse(string(tpl.Text))
	if err != nil {
		return nil
	}

	return t.Execute(writer, context)
}

func (e *engine) ExportInvoices(writer io.Writer, tpl *types.Template, invoices []*types.Invoice, account *types.Account) error {
	context := &types.InvoiceTemplateContext{
		Config:   e.cfg,
		Account:  account,
		Template: tpl,
		Invoices: invoices,
	}
	t, err := template.New(tpl.What).Parse(string(tpl.Text))
	if err != nil {
		return nil
	}

	return t.Execute(writer, context)
}
