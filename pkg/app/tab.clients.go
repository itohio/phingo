package app

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/phingo/pkg/types"
)

func (a *App) newClientsContents() fyne.CanvasObject {
	headers := []tableHeader{
		{"Name", 150},
		{"Description", 150},
		{"Invoice File Format", 150},
		{"ID", 150},
	}
	ret := a.newTable(
		func() int { return len(a.clients) },
		headers,
		func(row, col int, l *widget.Label) string {
			cl := a.clients[row]
			switch col {
			case 0:
				return cl.Name
			case 1:
				return cl.Description
			case 2:
				return cl.InvoiceFileNameFormat
			case 3:
				return cl.Id
			}
			return "-"
		},
		func(row int) {
			a.editClient(a.clients[row])
		},
	)

	return ret
}

func (a *App) editClient(cl *types.Client) {
	addContact := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {})

	eName := widget.NewEntry()
	eName.Text = cl.Name
	eDescription := widget.NewEntry()
	eDescription.Text = cl.Description

	eInvoiceFileFormat := widget.NewEntry()
	eInvoiceFileFormat.Text = cl.InvoiceFileNameFormat

	eNotes := widget.NewMultiLineEntry()
	eNotes.Text = strings.Join(cl.Notes, "\n")

	items := []*widget.FormItem{
		widget.NewFormItem("ID", widget.NewLabel(cl.Id)),
		widget.NewFormItem("Name", eName),
		widget.NewFormItem("Description", eDescription),
		widget.NewFormItem("Invoice File Format", eInvoiceFileFormat),
		widget.NewFormItem("Notes", eNotes),
		widget.NewFormItem("Contacts:", container.NewHBox(layout.NewSpacer(), addContact)),
	}
	form := widget.NewForm(items...)
	for k, v := range cl.Contact {
		addContactItem(form, k, v)
	}

	addContact.OnTapped = a.newContactItemAdder(form, cl.Contact)

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

func (a *App) newClient() {
	cl := &types.Client{
		Contact: a.makeDefaultContacts(),
	}
	a.editClient(cl)
}
