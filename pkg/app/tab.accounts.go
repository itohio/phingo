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

	eName := widget.NewEntry()
	eName.Text = acc.Name
	eDecimals := widget.NewEntry()
	eDecimals.Text = fmt.Sprint(acc.Decimals)
	eDecimals.Validator = validateIntMin(0)
	eDenomination := widget.NewEntry()
	eDenomination.Text = acc.Denom
	eInvoiceFileFormat := widget.NewEntry()
	eInvoiceFileFormat.Text = acc.InvoiceFileNameFormat
	eInvoiceCodeFormat := widget.NewEntry()
	eInvoiceCodeFormat.Text = acc.InvoiceCodeFormat
	eInvoiceDuePeriod := widget.NewEntry()
	eInvoiceDuePeriod.Text = fmt.Sprint(acc.InvoiceDuePeriod)
	eInvoiceDuePeriod.Validator = validateDurationMin(0)

	items := []*widget.FormItem{
		widget.NewFormItem("ID", widget.NewLabel(acc.Id)),
		widget.NewFormItem("Name", eName),
		widget.NewFormItem("Denomination", eDenomination),
		widget.NewFormItem("Decimals", eDecimals),
		widget.NewFormItem("Invoice File Format", eInvoiceFileFormat),
		widget.NewFormItem("Invoice Code Format", eInvoiceCodeFormat),
		widget.NewFormItem("Invoice Due Period", eInvoiceDuePeriod),
		widget.NewFormItem("Contacts:", container.NewHBox(layout.NewSpacer(), addContact)),
	}
	form := widget.NewForm(items...)
	for k, v := range acc.Contact {
		addContactItem(form, k, v)
	}

	addContact.OnTapped = a.newContactItemAdder(form, acc.Contact)

	dlg := dialog.NewCustomConfirm(
		"Account", "OK", "Cancel",
		container.NewScroll(form),
		func(b bool) {

		},
		a.window,
	)

	dlg.Resize(a.window.Canvas().Size())
	dlg.Show()
}

func (a *App) newAccount() {
	acc := &types.Account{
		Contact: a.makeDefaultContacts(),
	}
	a.editAccount(acc)
}
