package engine

import (
	"context"
	"image/color"
	"io"
	"os"

	"github.com/itohio/phingo/pkg/types"
	pdf "github.com/stephenafamo/goldmark-pdf"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

type engine struct {
	md goldmark.Markdown
}

func New(config *types.Config) (*engine, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.CJK,
			meta.New(
				meta.WithStoresInDocument(),
			),
		),
		goldmark.WithRenderer(
			pdf.New(
				pdf.WithTraceWriter(os.Stderr),
				pdf.WithContext(context.Background()),
				pdf.WithImageFS(os.DirFS(".")),
				pdf.WithLinkColor(color.RGBA{R: 0xCC, G: 0x45, B: 0x78, A: 255}),
				pdf.WithHeadingFont(pdf.GetTextFont("IBM Plex Serif", pdf.FontLora)),
				pdf.WithBodyFont(pdf.GetTextFont("Open Sans", pdf.FontRoboto)),
				pdf.WithCodeFont(pdf.GetCodeFont("Inconsolata", pdf.FontRobotoMono)),
			),
		),
	)

	return &engine{
		md: md,
	}, nil
}

func (e *engine) Meta(tpl *types.Template) (map[string]interface{}, error) {
	document := e.md.Parser().Parse(text.NewReader(tpl.Text))
	return document.OwnerDocument().Meta(), nil
}

func (e *engine) ExportProjects(writer io.Writer, tpl *types.Template, projects []*types.Project, account *types.Account) error {
	return e.md.Convert(tpl.Text, writer)
}

func (e *engine) ExportInvoices(writer io.Writer, tpl *types.Template, invoices []*types.Invoice, account *types.Account) error {
	return e.md.Convert(tpl.Text, writer)
}
