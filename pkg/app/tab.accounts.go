package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/phingo/pkg/types"
)

func (a *App) newAccountsContents() fyne.CanvasObject {
	headers := []tableHeader{
		{"Name", 150},
		{"Denomination", 150},
		{"Code Format", 150},
		{"File Format", 150},
		{"ID", 150},
	}
	ret := a.newTable(
		func() int { return len(a.accounts) },
		headers,
		func(row, col int, l *widget.Label) string {
			acc := a.accounts[row]
			switch col {
			case 0:
				return acc.Name
			case 1:
				return acc.Denom
			case 2:
				return acc.InvoiceCodeFormat
			case 3:
				return acc.InvoiceFileNameFormat
			case 4:
				return acc.Id
			}
			return "-"
		},
		func(row int) {
			a.editAccount(a.accounts[row])
		},
	)

	return ret
}

func (a *App) editAccount(acc *types.Account) {
	addContact := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {})
	items := []*widget.FormItem{
		widget.NewFormItem("ID", widget.NewLabel(acc.Id)),
		widget.NewFormItem("Name", widget.NewLabel(acc.Name)),
		widget.NewFormItem("Denomination", widget.NewLabel(acc.Denom)),
		widget.NewFormItem("Decimals", widget.NewLabel(fmt.Sprint(acc.Decimals))),
		widget.NewFormItem("Invoice File Format", widget.NewLabel(acc.InvoiceFileNameFormat)),
		widget.NewFormItem("Invoice Code Format", widget.NewLabel(acc.InvoiceCodeFormat)),
		widget.NewFormItem("Invoice Due Period", widget.NewLabel(fmt.Sprint(acc.InvoiceDuePeriod))),
		widget.NewFormItem("Contacts:", container.NewHBox(layout.NewSpacer(), addContact)),
	}
	addContact.OnTapped = func() {
		eType := widget.NewEntry()
		eValue := widget.NewEntry()
		dialog.ShowForm(
			"Contact", "OK", "Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("Type", eType),
				widget.NewFormItem("Value", eValue),
			},
			func(b bool) {
				if !b {
					return
				}
			},
			a.window,
		)
	}

	for k, v := range acc.Contact {
		btn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {})
		fi := widget.NewFormItem(
			k,
			container.NewHBox(widget.NewLabel(v), layout.NewSpacer(), btn),
		)
		items = append(items, fi)
	}

	dlg := dialog.NewForm(
		"Account", "OK", "Cancel",
		items,
		func(b bool) {

		},
		a.window,
	)

	dlg.Resize(a.window.Canvas().Size())
	dlg.Show()
}
