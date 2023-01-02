package engine

import (
	"errors"
	"io"
	"io/fs"

	html "github.com/itohio/phingo/pkg/engine/html"
	markdown "github.com/itohio/phingo/pkg/engine/markdown"
	pdf "github.com/itohio/phingo/pkg/engine/pdf"
	"github.com/itohio/phingo/pkg/types"
)

// Engine interface represents various exporting engines
type Engine interface {
	// Meta will extract meta information from a template
	Meta(tpl *types.Template) (map[string]interface{}, error)

	ExportAccounts(writer io.Writer, tpl *types.Template, accounts []*types.Account) error
	ExportClients(writer io.Writer, tpl *types.Template, clients []*types.Client) error
	ExportProjects(writer io.Writer, tpl *types.Template, projects []*types.Project) error
	ExportInvoices(writer io.Writer, tpl *types.Template, invoices []*types.Invoice, account *types.Account) error
}

// New will attempt to create an exporting engine given the name of the engine,
// general configuration and a file system that contains assets that
// may be referenced in the templates.
func New(what string, config *types.Config, fsys fs.FS) (Engine, error) {
	switch what {
	case "pdf":
		return pdf.New(config, fsys)
	case "htm":
		fallthrough
	case "html":
		return html.New(config, fsys)
	case "console":
		fallthrough
	case "markdown":
		return markdown.New(config)
	default:
		return nil, errors.New("export engine unknown")
	}
}
