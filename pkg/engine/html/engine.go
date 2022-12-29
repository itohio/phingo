package engine

import (
	"bytes"
	"io"

	markdown "github.com/itohio/phingo/pkg/engine/markdown"
	"github.com/itohio/phingo/pkg/types"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type Engine struct {
	mde *markdown.Engine
	md  goldmark.Markdown
}

func New(config *types.Config) (*Engine, error) {
	mde, err := markdown.New(config)
	if err != nil {
		return nil, err
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.CJK,
			meta.New(
				meta.WithStoresInDocument(),
			),
			extension.NewTable(
				extension.WithTableCellAlignMethod(extension.TableCellAlignAttribute),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithHeadingAttribute(),
			parser.WithASTTransformers(util.Prioritized(&TableASTTransformer{}, 1000)),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	return &Engine{
		md:  md,
		mde: mde,
	}, nil
}

func (e *Engine) Meta(tpl *types.Template) (map[string]interface{}, error) {
	document := e.md.Parser().Parse(text.NewReader(tpl.Text))
	return document.OwnerDocument().Meta(), nil
}

func (e *Engine) ExportAccounts(writer io.Writer, tpl *types.Template, accounts []*types.Account) error {
	buf := bytes.NewBuffer(nil)
	if err := e.mde.ExportAccounts(buf, tpl, accounts); err != nil {
		return err
	}
	return e.md.Convert(buf.Bytes(), writer)
}

func (e *Engine) ExportClients(writer io.Writer, tpl *types.Template, clients []*types.Client) error {
	buf := bytes.NewBuffer(nil)
	if err := e.mde.ExportClients(buf, tpl, clients); err != nil {
		return err
	}
	return e.md.Convert(buf.Bytes(), writer)
}

func (e *Engine) ExportProjects(writer io.Writer, tpl *types.Template, projects []*types.Project) error {
	buf := bytes.NewBuffer(nil)
	if err := e.mde.ExportProjects(buf, tpl, projects); err != nil {
		return err
	}
	return e.md.Convert(buf.Bytes(), writer)
}

func (e *Engine) ExportInvoices(writer io.Writer, tpl *types.Template, invoices []*types.Invoice, account *types.Account) error {
	buf := bytes.NewBuffer(nil)
	if err := e.mde.ExportInvoices(buf, tpl, invoices, account); err != nil {
		return err
	}
	return e.md.Convert(buf.Bytes(), writer)
}
