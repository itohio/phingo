package engine

import (
	"errors"
	"io"

	console "github.com/itohio/phingo/pkg/engine/console"
	pdf "github.com/itohio/phingo/pkg/engine/pdf"
	"github.com/itohio/phingo/pkg/types"
)

type Engine interface {
	// Meta will extract meta information from a template
	Meta(tpl *types.Template) (map[string]interface{}, error)

	ExportProjects(writer io.Writer, tpl *types.Template, projects []*types.Project, account *types.Account) error
	ExportInvoices(writer io.Writer, tpl *types.Template, invoices []*types.Invoice, account *types.Account) error
}

func New(what string, config *types.Config) (Engine, error) {
	switch what {
	case "pdf":
		return pdf.New(config)
	case "console":
		return console.New(config)
	default:
		return nil, errors.New("export engine unknown")
	}
}
