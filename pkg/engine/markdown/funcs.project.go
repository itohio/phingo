package engine

import (
	"text/template"
	"time"

	"github.com/itohio/phingo/pkg/types"
)

func makeProjectFuncs(context *types.ProjectTemplateContext) template.FuncMap {
	return addDefaultFuncs(template.FuncMap{
		"TotalProgress": func(p *types.Project) float32 {
			if p == nil {
				return 0
			}
			return p.TotalProgress()
		},
		"TotalPriceWords": func(p *types.Project) string {
			pr := p.TotalPrice()
			if pr == nil {
				return "-"
			}
			return pr.Words()
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
				return context.Config.Format(pr)
			}
			pr := l.Price(p)
			if pr == nil {
				return "-"
			}
			return context.Config.Format(pr)
		},
		"Duration": func(l *types.Project_LogEntry) string {
			if l == nil {
				return "-"
			}
			return time.Duration(l.Duration).String()
		},
		"Contacts": makeContactsFunc(context.Config),
		"Account":  makeAccountFunc(context.Config),
		"Client":   makeClientFunc(context.Config),
	})
}
