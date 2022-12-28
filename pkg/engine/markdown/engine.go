package engine

import (
	"io"
	"text/template"
	"time"

	"github.com/itohio/phingo/pkg/types"
)

type Engine struct {
	cfg *types.Config
}

func New(config *types.Config) (*Engine, error) {
	return &Engine{
		cfg: config,
	}, nil
}

func (e *Engine) Meta(tpl *types.Template) (map[string]interface{}, error) {
	return nil, nil
}

func (e *Engine) ExportAccounts(writer io.Writer, tpl *types.Template, accounts []*types.Account) error {
	context := &types.AccountTemplateContext{
		Config:   e.cfg,
		Template: tpl,
		Accounts: accounts,
	}
	t, err := template.New(tpl.What).Parse(string(tpl.Text))
	if err != nil {
		return err
	}

	return t.Execute(writer, context)
}

func (e *Engine) ExportClients(writer io.Writer, tpl *types.Template, clients []*types.Client) error {
	context := &types.ClientTemplateContext{
		Config:   e.cfg,
		Template: tpl,
		Clients:  clients,
	}
	t, err := template.New(tpl.What).Parse(string(tpl.Text))
	if err != nil {
		return err
	}

	return t.Execute(writer, context)
}

func (e *Engine) ExportProjects(writer io.Writer, tpl *types.Template, projects []*types.Project) error {
	context := &types.ProjectTemplateContext{
		Config:   e.cfg,
		Template: tpl,
		Projects: projects,
	}
	t := template.New(tpl.What)
	t.Funcs(template.FuncMap{
		"TotalProgress": func(p *types.Project) float32 {
			if p == nil {
				return 0
			}
			return p.TotalProgress()
		},
		"TotalPrice": func(p *types.Project) string {
			pr := p.TotalPrice()
			if pr == nil {
				return "-"
			}
			return pr.Pretty()
		},
		"TotalDuration": func(p *types.Project) string {
			if p == nil {
				return "-"
			}
			return p.TotalDuration().String()
		},
		"Rate": func(p *types.Project) string {
			if p == nil {
				return "-"
			}
			return p.RateString()
		},
		"Price": func(p *types.Project, l *types.Project_LogEntry) string {
			if p == nil {
				return "-"
			}
			if l == nil {
				return "-"
			}
			pr := l.Price(p)
			return pr.Pretty()
		},
		"Duration": func(l *types.Project_LogEntry) string {
			if l == nil {
				return "-"
			}
			return time.Duration(l.Duration).String()
		},
	})
	t, err := t.Parse(string(tpl.Text))
	if err != nil {
		return err
	}

	return t.Execute(writer, context)
}

func (e *Engine) ExportInvoices(writer io.Writer, tpl *types.Template, invoices []*types.Invoice, account *types.Account) error {
	context := &types.InvoiceTemplateContext{
		Config:   e.cfg,
		Account:  account,
		Template: tpl,
		Invoices: invoices,
	}
	t, err := template.New(tpl.What).Parse(string(tpl.Text))
	if err != nil {
		return err
	}

	return t.Execute(writer, context)
}
