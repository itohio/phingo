package engine

import (
	"text/template"
	"time"

	"github.com/itohio/phingo/pkg/types"
)

func makeInvoiceFuncs(context *types.InvoiceTemplateContext) template.FuncMap {
	return addDefaultFuncs(template.FuncMap{
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
				pr := p.TotalPrice()
				if pr == nil {
					return "-"
				}
				return pr.Pretty()
			}
			pr := l.Price(p)
			if pr == nil {
				return "-"
			}
			return pr.Pretty()
		},
		"Duration": func(l *types.Project_LogEntry) string {
			if l == nil {
				return "-"
			}
			return time.Duration(l.Duration).String()
		},
		"Client":   makeClientFunc(context.Config),
		"Account":  makeAccountFunc(context.Config),
		"Contacts": makeContactsFunc(context.Config),
	})
}
