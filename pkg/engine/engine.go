package engine

import (
	"errors"
	"io"

	html "github.com/itohio/phingo/pkg/engine/html"
	markdown "github.com/itohio/phingo/pkg/engine/markdown"
	pdf "github.com/itohio/phingo/pkg/engine/pdf"
	"github.com/itohio/phingo/pkg/types"
)

type Engine interface {
	// Meta will extract meta information from a template
	Meta(tpl *types.Template) (map[string]interface{}, error)

	ExportAccounts(writer io.Writer, tpl *types.Template, accounts []*types.Account) error
	ExportClients(writer io.Writer, tpl *types.Template, clients []*types.Client) error
	ExportProjects(writer io.Writer, tpl *types.Template, projects []*types.Project) error
	ExportInvoices(writer io.Writer, tpl *types.Template, invoices []*types.Invoice, account *types.Account) error
}

func New(what string, config *types.Config) (Engine, error) {
	switch what {
	case "pdf":
		return pdf.New(config)
	case "htm":
		fallthrough
	case "html":
		return html.New(config)
	case "console":
		fallthrough
	case "markdown":
		return markdown.New(config)
	default:
		return nil, errors.New("export engine unknown")
	}
}
