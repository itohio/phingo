package engine

import (
	gast "github.com/yuin/goldmark/ast"
	tast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type TableASTTransformer struct {
}

func (a *TableASTTransformer) Transform(node *gast.Document, reader text.Reader, pc parser.Context) {
	for c := node.FirstChild(); c != nil; {
		switch c.Kind() {
		case tast.KindTable:
			a.transformTable(c)
		case gast.KindDocument:
			a.transformDocument(c)
		}
		c = c.NextSibling()
	}
}

func (a *TableASTTransformer) transformTable(node gast.Node) {
	node.SetAttributeString("width", []byte("100%"))
	node.SetAttributeString("border", []byte("1px solid grey"))
	node.SetAttributeString("bgcolor", []byte("#F9F9F9"))
	node.SetAttributeString("cellspacing", []byte("0"))
	node.SetAttributeString("cellpadding", []byte("5px"))
	row := 0
	for c := node.FirstChild(); c != nil; {
		switch c.Kind() {
		case tast.KindTableHeader:
			c.SetAttributeString("bgcolor", []byte("#BBB"))
		case tast.KindTableRow:
			if row%2 == 0 {
				c.SetAttributeString("bgcolor", []byte("#f0f0f0"))
			}
			row++
		}
		c = c.NextSibling()
	}
}

func (a *TableASTTransformer) transformDocument(node gast.Node) {
}
