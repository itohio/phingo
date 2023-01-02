package engine

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"image/color"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	markdown "github.com/itohio/phingo/pkg/engine/markdown"
	"github.com/itohio/phingo/pkg/types"
	fences "github.com/stefanfritsch/goldmark-fences"
	pdf "github.com/stephenafamo/goldmark-pdf"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Engine struct {
	mde *markdown.Engine
	fs  fs.FS
}

func New(config *types.Config, fsys fs.FS) (*Engine, error) {
	mde, err := markdown.New(config)
	if err != nil {
		return nil, err
	}

	ret := &Engine{
		mde: mde,
		fs:  fsys,
	}

	return ret, nil
}

func getValue[T any](m map[string]interface{}, key string, fallback T) T {
	if m == nil {
		return fallback
	}
	if val, ok := m[key]; ok {
		if ret, ok := val.(T); ok {
			return ret
		}
	}
	return fallback
}

func getColor(m map[string]interface{}, key string, fallback color.Color) color.Color {
	clr := getValue(m, key, "")
	if clr == "" {
		return fallback
	}

	clr = strings.TrimPrefix(clr, "#")

	b, err := hex.DecodeString(clr)
	if err != nil || len(b) < 1 {
		return fallback
	}

	switch len(b) {
	case 1:
		return color.RGBA{b[0], b[0], b[0], 255}
	case 3:
		return color.RGBA{b[0], b[1], b[2], 255}
	case 4:
		return color.RGBA{b[0], b[1], b[2], b[3]}
	}
	return fallback
}

func (e *Engine) makeRenderer(tpl *types.Template, title, author string, creationTime time.Time) goldmark.Markdown {
	cfg, err := e.Meta(tpl)
	if err != nil {
		log.Println("Could not extract template config: ", err.Error())
	}
	log.Println(cfg)

	options := []pdf.Option{
		pdf.WithFpdf(
			context.Background(),
			pdf.FpdfConfig{
				Subject:      strings.ToTitle(tpl.What),
				Title:        title,
				Author:       author,
				Orientation:  getValue(cfg, "orientation", ""),
				PaperSize:    getValue(cfg, "paperSize", "a4"),
				Creator:      "phingo.itohio",
				CreationTime: creationTime,
			},
		),
		pdf.WithContext(context.Background()),
		pdf.WithImageFS(e.fs),
		pdf.WithLinkColor(getColor(cfg, "linkColor", color.RGBA{R: 0xCC, G: 0x45, B: 0x78, A: 255})),
		pdf.WithHeadingFont(pdf.GetTextFont(getValue(cfg, "headingFont", "IBM Plex Serif"), pdf.FontLora)),
		pdf.WithBodyFont(pdf.GetTextFont(getValue(cfg, "bodyFont", "Open Sans"), pdf.FontRoboto)),
		pdf.WithCodeFont(pdf.GetCodeFont(getValue(cfg, "codeFont", "Inconsolata"), pdf.FontRobotoMono)),
	}

	if getValue(cfg, "trace", false) {
		options = append(options, pdf.WithTraceWriter(os.Stderr))
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.CJK,
			meta.New(
				meta.WithStoresInDocument(),
			),
			&fences.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithHeadingAttribute(),
			// parser.WithASTTransformers(util.Prioritized(&TableASTTransformer{}, 1000)),
		),
		goldmark.WithRenderer(
			pdf.New(options...),
		),
	)
	return md
}

func (e *Engine) Meta(tpl *types.Template) (map[string]interface{}, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.New(
				meta.WithStoresInDocument(),
			),
		),
	)
	document := md.Parser().Parse(text.NewReader(tpl.Text))
	return document.OwnerDocument().Meta(), nil
}

func (e *Engine) ExportAccounts(writer io.Writer, tpl *types.Template, accounts []*types.Account) error {
	buf := bytes.NewBuffer(nil)
	if err := e.mde.ExportAccounts(buf, tpl, accounts); err != nil {
		return err
	}
	var (
		author string
		tm     time.Time = time.Now()
	)
	if len(accounts) > 0 {
		var ok bool
		if author, ok = accounts[0].Contact[types.ContactFullName]; !ok {
			if author, ok = accounts[0].Contact[types.ContactName]; !ok {
				author = accounts[0].Name
			}
		}
	}
	md := e.makeRenderer(tpl, "Accounts", author, tm)
	return md.Convert(buf.Bytes(), writer)
}
func (e *Engine) ExportClients(writer io.Writer, tpl *types.Template, clients []*types.Client) error {
	buf := bytes.NewBuffer(nil)
	if err := e.mde.ExportClients(buf, tpl, clients); err != nil {
		return err
	}
	md := e.makeRenderer(tpl, "Clients", "Phingo", time.Now())
	return md.Convert(buf.Bytes(), writer)
}

func (e *Engine) ExportProjects(writer io.Writer, tpl *types.Template, projects []*types.Project) error {
	buf := bytes.NewBuffer(nil)
	if err := e.mde.ExportProjects(buf, tpl, projects); err != nil {
		return err
	}
	var (
		author string
		tm     time.Time = time.Now()
	)
	if len(projects) > 0 && projects[0].Account != nil {
		var ok bool
		if author, ok = projects[0].Account.Contact[types.ContactFullName]; !ok {
			if author, ok = projects[0].Account.Contact[types.ContactName]; !ok {
				author = projects[0].Account.Name
			}
		}
	}
	md := e.makeRenderer(tpl, "Projects", author, tm)
	return md.Convert(buf.Bytes(), writer)
}

func (e *Engine) ExportInvoices(writer io.Writer, tpl *types.Template, invoices []*types.Invoice, account *types.Account) error {
	buf := bytes.NewBuffer(nil)
	if err := e.mde.ExportInvoices(buf, tpl, invoices, account); err != nil {
		return err
	}
	var (
		author string
		title  string
		tm     time.Time = time.Now()
	)
	if len(invoices) > 0 {
		acc := invoices[0].Account
		if account != nil {
			acc = account
		}
		var ok bool
		if author, ok = acc.Contact[types.ContactFullName]; !ok {
			if author, ok = acc.Contact[types.ContactName]; !ok {
				author = acc.Name
			}
		}
		title = fmt.Sprint("Invoice", invoices[0].Id)
	}
	md := e.makeRenderer(tpl, title, author, tm)
	return md.Convert(buf.Bytes(), writer)
}
