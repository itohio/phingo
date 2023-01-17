package engine

import (
	"text/template"

	"github.com/itohio/phingo/pkg/bi"
	"github.com/itohio/phingo/pkg/types"
)

func makeInvoiceFuncs(context *types.InvoiceTemplateContext) template.FuncMap {
	return addDefaultFuncs(template.FuncMap{
		"Summary": func(inv *types.Invoice) *bi.InvoiceSummary {
			if inv == nil {
				return nil
			}
			return bi.NewInvoiceSummary(inv, context.SelectedAccount(inv).Denom, uint32(context.SelectedAccount(inv).Decimals))
		},
		"ItemPrice": func(inv *types.Invoice, item *types.Invoice_Item) *types.Price {
			if item == nil {
				return nil
			}
			return item.Price(context.SelectedAccount(inv).Denom)
		},
		"Client":   makeClientFunc(context.Config),
		"Account":  makeAccountFunc(context.Config),
		"Contacts": makeContactsFunc(context.Config),
		"Convert": func(a, b *types.Price) *types.Price {
			return a.Convert(b)
		},
	})
}
