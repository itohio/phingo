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

func (e *engine) ExportProjects(writer io.Writer, tpl *types.Template, projects []*types.Project, account *types.Account) error {
	context := &types.ProjectTemplateContext{
		Config:   e.cfg,
		Account:  account,
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
