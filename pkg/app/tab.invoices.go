package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/phingo/pkg/bi"
	"github.com/itohio/phingo/pkg/types"
)

func (a *App) newInvoicesContents() fyne.CanvasObject {
	headers := []tableHeader{
		{"Code", 150},
		{"Date", 150},
		{"Project", 150},
		{"Client", 150},
		{"Account", 150},
		{"Subtotal", 150},
		{"Tax", 150},
		{"Total", 150},
		{"Payed", 150},
	}
	ret := a.newTable(
		func() int { return len(a.invoices) },
		headers,
		func(row, col int, l *widget.Label) string {
			inv := a.invoices[row]
			switch col {
			case 0:
				return inv.Code
			case 1:
				return inv.IssueDate
			case 2:
				return inv.Project.Name
			case 3:
				return inv.Client.Name
			case 4:
				return inv.Account.Name
			case 5:
				summary := bi.NewInvoiceSummary(inv, inv.Account.Denom, uint32(inv.Account.Decimals))
				return summary.Subtotal.Pretty()
			case 6:
				summary := bi.NewInvoiceSummary(inv, inv.Account.Denom, uint32(inv.Account.Decimals))
				return summary.Tax.Pretty()
			case 7:
				summary := bi.NewInvoiceSummary(inv, inv.Account.Denom, uint32(inv.Account.Decimals))
				return summary.Total.Pretty()
			case 8:
				return "unpayed"
			}
			return "-"
		},
		nil,
	)

	return ret
}

func (a *App) editInvoice(inv *types.Invoice) {
}

func (a *App) newInvoice() {
	inv := &types.Invoice{}
	a.editInvoice(inv)
}
